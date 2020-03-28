package models

type PropertyName struct {
	BaseModel

	Name           string          `gorm:"type: varchar(50); not null" json:"name"`
	PropertyValues []PropertyValue `json:"property_values"`
}

func (pname *PropertyName) Find() (err error) {
	err = db.First(&pname).Error
	return
}
