package models

import (
	"github.com/jinzhu/gorm"
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

func (delivery Delivery) DestroyRules() {
	db.Where("delivery_id = ?", delivery.ID).Delete(DeliveryRule{})
}

func (delivery *Delivery) Reload() (err error) {
	err = Find(delivery, Options{})
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
