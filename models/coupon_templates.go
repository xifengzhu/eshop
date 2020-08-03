package models

import (
	"encoding/json"
	"github.com/captaincodeman/couponcode"
	"log"
	"time"
)

type CouponTemplate struct {
	BaseModel

	Code         string     `gorm:"type: varchar(20);unique_index" json:"code"`
	Name         string     `gorm:"type: varchar(64); not null" json:"name"`
	Kind         string     `gorm:"type: varchar(64);" json:"kind"`
	Creator      string     `gorm:"type: varchar(64);" json:"creator"`
	Stock        int        `gorm:"type: int;" json:"stock"`
	CatchLimit   int        `gorm:"type: int;" json:"catch_limit"`
	StartAt      *time.Time `gorm:"type: datetime; default: CURRENT_TIMESTAMP" json:"start_at,omitempty"`
	EndAt        *time.Time `gorm:"type: datetime;" json:"end_at,omitempty"`
	Configs      JSON       `gorm:"type: json;" json:"configs"`
	Coupons      []*Coupon  `json:"coupons,omitempty"`
	CouponsCount int        `gorm:"type: int; default:0" json:"coupons_count"`
}

func (template *CouponTemplate) BeforeCreate() (err error) {
	var code string
	exist := true
	for exist {
		code = couponcode.Generate()
		parmMap := map[string]interface{}{"code": code}
		exist = Exist(&CouponTemplate{}, Options{Conditions: parmMap})
	}
	template.Code = code
	return
}

func (template *CouponTemplate) CatchCountForUser(userID int) (count int) {
	db.Model(&Coupon{}).Where("coupon_template_id = ? AND user_id = ?", template.ID, userID).Count(&count)
	return
}

func (template *CouponTemplate) GenerateCoupon(qty int) {
	coupons := template.GenerateCouponsData(qty)
	for _, coupon := range coupons {
		db.Create(&coupon)
	}
}

func (template *CouponTemplate) GenerateCouponsData(qty int) (coupons []Coupon) {
	var startAt time.Time
	var endAt time.Time
	configs := template.ConfigJSON()
	if template.isFixTerm() {
		startAt = time.Now()
		endAt = startAt.AddDate(0, 0, int(configs["fix_term"].(float64)))
	} else {
		startAt, _ = time.Parse(time.RFC3339, configs["start_at"].(string))
		endAt, _ = time.Parse(time.RFC3339, configs["end_at"].(string))
	}

	for i := 0; i < qty; i++ {
		coupon := Coupon{
			CouponTemplateID: template.ID,
			Name:             template.Name,
			Kind:             template.Kind,
			CatchLimit:       template.CatchLimit,
			StartAt:          &startAt,
			EndAt:            &endAt,
			Configs:          template.Configs,
		}
		coupons = append(coupons, coupon)
	}
	return coupons
}

func (template *CouponTemplate) isFixTerm() bool {
	return template.ConfigJSON()["date_type"] == "fix_term"
}

func (template *CouponTemplate) ConfigJSON() map[string]interface{} {
	var data map[string]interface{}

	if err := json.Unmarshal(template.Configs, &data); err != nil {
		log.Println("===template config err===", err)
	}
	return data
}
