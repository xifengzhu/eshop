package models

import (
	// "fmt"
	"github.com/jinzhu/gorm"
	"github.com/xifengzhu/eshop/helpers/utils"
)

type Delivery struct {
	BaseModel

	WxappId string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name    string `gorm:"type: varchar(255); not null" json:"name"`
	// 1 为按件计费 2 按重量计费
	Way           int            `gorm:"type: tinyint; " json:"way"`
	DeliveryRules []DeliveryRule `json:"delivery_rules,omitempty"`
}

func (Delivery) TableName() string {
	return "delivery"
}

func (d *Delivery) IsCaculateByNum() bool {
	return d.Way == 1
}

func (d *Delivery) IsCaculateByWeight() bool {
	return d.Way == 2
}

func (cat Delivery) All(pagination *utils.Pagination, maps map[string]interface{}) (categories []Delivery, err error) {
	offset := (pagination.Page - 1) * pagination.PerPage
	db.Where(maps).Offset(offset).Limit(pagination.PerPage).Order("position asc").Find(&categories)
	pagination.Total = GetDeliveryTotal(maps)
	return
}

func (cat Delivery) AllWithoutPagination() (categories []Delivery, err error) {
	err = db.Preload("Children").Order("position asc").Where("parent_id IS NULL").Find(&categories).Error
	return
}

func GetDeliveryTotal(maps interface{}) (count int) {
	db.Model(&Delivery{}).Where(maps).Count(&count)
	return
}

func (delivery *Delivery) Create() (err error) {
	err = db.Create(&delivery).Error
	return
}

func (delivery Delivery) DestroyRules() {
	db.Where("delivery_id = ?", delivery.ID).Delete(DeliveryRule{})
}

func (delivery *Delivery) Find() (err error) {
	err = db.Preload("DeliveryRules").Find(&delivery).Error
	return
}

func (delivery *Delivery) Reload() (err error) {
	err = delivery.Find()
	return
}

func (delivery *Delivery) Save() (err error) {
	err = db.Save(&delivery).Error
	return
}

func (delivery *Delivery) NestUpdate() (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var rules []DeliveryRule
		var deleteIDs []int
		for _, deliveryRule := range delivery.DeliveryRules {
			if deliveryRule.Destroy == true {
				deleteIDs = append(deleteIDs, deliveryRule.ID)
				// tx.Delete(&deliveryRule)
			} else {
				rules = append(rules, deliveryRule)
			}
		}
		if len(deleteIDs) > 0 {
			err = tx.Where("id = ?", deleteIDs).Delete(DeliveryRule{}).Error
		}
		delivery.DeliveryRules = rules
		err = tx.Save(&delivery).Error
		return err
	})
}
