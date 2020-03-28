package models

type PropertyValue struct {
	BaseModel

	Name           string `gorm:"type: varchar(50); not null" json:"name"`
	PropertyNameID int   `gorm:"type: int; not null" json:"property_name_id"`
}

func (pvalue *PropertyValue) Find() (err error) {
	err = db.First(&pvalue).Error
	return
}
