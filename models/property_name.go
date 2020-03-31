package models

import (
	"github.com/jinzhu/gorm"
)

type PropertyName struct {
	BaseModel

	Name           string          `gorm:"type: varchar(50); not null" json:"name"`
	PropertyValues []PropertyValue `json:"property_values"`
}

func (PropertyName) TableName() string {
	return "property_name"
}

func (p *PropertyName) RemovePropertyValues() {
	db.Where("property_name_id = ?", p.ID).Delete(PropertyValue{})
}

func (pname *PropertyName) NestUpdate() (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var values []PropertyValue
		var deleteIDs []int
		for _, goods := range pname.PropertyValues {
			if goods.Destroy == true {
				deleteIDs = append(deleteIDs, goods.ID)
			} else {
				values = append(values, goods)
			}
		}
		if len(deleteIDs) > 0 {
			err = tx.Where("id = ?", deleteIDs).Delete(PropertyValue{}).Error
		}
		pname.PropertyValues = values
		err = tx.Save(&pname).Error
		return err
	})
}
