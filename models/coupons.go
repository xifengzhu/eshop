package models

import (
	"encoding/json"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/xifengzhu/eshop/helpers/utils"
	config "github.com/xifengzhu/eshop/initializers"
	// "github.com/xifengzhu/eshop/models/coupons"
	"log"
	"time"
)

// type Config struct {
// 	MinAmount    float64    `json:"min_amount" validate:"required,gte=0"`
// 	ResourceType string     `json:"resource_type,omitempty"`
// 	Resources    []int      `json:"resources,omitempty" validate:"required_with=ResourceType"`
// 	DateType     string     `json:"date_type" validate:"required,oneof='fix_term' 'time_range'"`
// 	FixTerm      int        `json:"fix_term,omitempty" validate:"rfe=DateType:fix_term"`
// 	StartAt      *time.Time `json:"start_at,omitempty" validate:"required_with=EndAt,rfe=DateType:time_range"`
// 	EndAt        *time.Time `json:"end_at,omitempty" validate:"required_with=StartAt,rfe=DateType:time_range""`

// 	ReduceAmount float64 `json:"reduce_amount,omitempty" validate:"required_without=Percentage"`
// 	Percentage   int     `json:"percentage,omitempty" validate:"required_without=ReduceAmount"`
// }

type Coupon struct {
	BaseModel

	Code             string          `gorm:"type: varchar(20);unique_index" json:"code"`
	Name             string          `gorm:"type: varchar(64); not null" json:"name"`
	Kind             string          `gorm:"type: varchar(64);" json:"kind"`
	State            string          `gorm:"type: enum('actived', 'lock', 'used', 'expired');default:'actived'" json:"state"`
	CatchLimit       int             `gorm:"type: int;" json:"catch_limit"`
	StartAt          *time.Time      `gorm:"type: datetime;" json:"start_at,omitempty"`
	EndAt            *time.Time      `gorm:"type: datetime;" json:"end_at,omitempty"`
	Configs          JSON            `gorm:"type: json;" json:"configs"`
	UserID           int             `gorm:"type: int;" json:"user_id"`
	CouponTemplateID int             `gorm:"type: int;" json:"coupon_template_id"`
	User             *User           `json:"user,omitempty"`
	CouponTemplate   *CouponTemplate `json:"coupon_template,omitempty"`
	Order            *CouponOrder    `gorm:"-" json:"coupon_order,omitempty"`
	Adjustment       Adjustment      `gorm:"polymorphic:Source;" json:"adjustment,omitempty"`
}

func (coupon *Coupon) BeforeCreate() (err error) {
	var code string
	exist := true
	for exist {
		code = utils.RandStringRunes(6)
		parmMap := map[string]interface{}{"code": code}
		exist = Exist(&Coupon{}, Options{Conditions: parmMap})
	}
	coupon.Code = code
	return
}

func (coupon *Coupon) updateTemplateCounterCache(tx *gorm.DB) (err error) {
	tx.Model(&CouponTemplate{}).Where("id = ?", coupon.CouponTemplateID).Update("coupons_count", gorm.Expr("coupons_count + ?", 1))
	return
}

func (coupon *Coupon) AfterCreate(tx *gorm.DB) (err error) {
	log.Println("========coupon after create======")
	coupon.updateTemplateCounterCache(tx)
	coupon.enqueueExpireCouponJob()
	return
}

func (c *Coupon) Expire() {
	Update(c, map[string]interface{}{"state": "expired"})
}

func (coupon *Coupon) enqueueExpireCouponJob() {
	log.Println("========coupon enqueueExpireCouponJob======")
	duration := (*coupon.EndAt).Sub(*coupon.StartAt).Seconds()
	config.JobEnqueuer.EnqueueIn("expire_coupon", int64(duration), work.Q{"coupon_id": coupon.ID})
}

func (coupon *Coupon) ConfigJSON() map[string]interface{} {
	var data map[string]interface{}

	if err := json.Unmarshal(coupon.Configs, &data); err != nil {
		log.Println("===coupon config err===", err)
	}
	return data
}

// =============  coupon logics ==============

type CouponOrderItem struct {
	ProductAmount float64 `json:"product_amount"`
	ReduceAmount  float64 `json:"reduce_amount"`
	ResourceID    int     `json:"resource_id"`
	ResourceType  string  `json:"resource_type"`
}

type CouponOrder struct {
	ProductAmount float64           `json:"product_amount"`
	FreightFee    float64           `json:"freight_fee"`
	ReduceAmount  float64           `json:"reduce_amount"`
	ReduceFreight float64           `json:"reduce_freight"`
	Items         []CouponOrderItem `json:"items"`
}

func (coupon Coupon) Apply(order CouponOrder) CouponOrder {
	coupon.Order = &order
	coupon.caculateItemsReduceAmount()
	coupon.Order.FreightFee = 0
	return *coupon.Order
}

func (coupon Coupon) IsEligible() bool {
	return coupon.ConfigJSON()["min_amount"].(float64) <= coupon.couponAmountInOrder()
}

func (coupon Coupon) isCouponType() bool {
	return len(coupon.avaliableResources()) > 0
}

func (coupon Coupon) couponAmountInOrder() (amount float64) {
	items := coupon.effectItems()
	for _, item := range items {
		amount += item.ProductAmount
	}
	return
}

func (coupon Coupon) effectItems() (items []CouponOrderItem) {
	if coupon.isCouponType() {
		for _, item := range coupon.Order.Items {
			if coupon.includeResource(item.ResourceID) {
				items = append(items, item)
			}
		}
	} else {
		items = coupon.Order.Items
	}
	return
}

// resource_type: product, category
func (coupon Coupon) resourceType() string {
	return coupon.ConfigJSON()["resource_type"].(string)
}

func (coupon Coupon) resourceIDs() (ids []int) {
	if coupon.ConfigJSON()["resource_ids"] != nil {
		ids = coupon.ConfigJSON()["resource_ids"].([]int)
	}
	return
}

func (coupon Coupon) avaliableResources() (products []Product) {
	if len(coupon.resourceIDs()) > 0 {
		if coupon.resourceType() == "product" {
			db.Model(&Product{}).Where("id IN (?)", coupon.resourceIDs()).Find(&products)
		} else {
			db.Model(&Category{}).Where("id IN (?)", coupon.resourceIDs()).Related(&products)
		}
	}
	return
}

func (coupon Coupon) includeResource(resourceID int) (result bool) {
	resources := coupon.avaliableResources()
	if len(resources) > 0 {
		for _, item := range resources {
			if item.ID == resourceID {
				result = true
				return
			}
		}
		return
	} else {
		result = true
		return
	}
}

func (coupon Coupon) reduceAmount() float64 {
	if coupon.Kind == "fixed_amount" {
		return coupon.ConfigJSON()["reduce_amount"].(float64)
	} else {
		return coupon.couponAmountInOrder() * (1 - coupon.ConfigJSON()["percentage"].(float64))
	}
}

func (coupon *Coupon) caculateItemsReduceAmount() {
	if coupon.IsEligible() {
		reduceAmount := coupon.ConfigJSON()["reduce_amount"].(float64)
		coupon.Order.ReduceAmount = reduceAmount
		for i, item := range coupon.Order.Items {
			var itemReduceAmount float64
			if coupon.includeResource(item.ResourceID) {
				itemReduceAmount = item.ProductAmount / coupon.couponAmountInOrder() * coupon.reduceAmount()
			}
			coupon.Order.Items[i].ReduceAmount = itemReduceAmount
		}
	} else {
		for i, _ := range coupon.Order.Items {
			coupon.Order.Items[i].ReduceAmount = 0
		}
	}
}
