package models

type PropertyName struct {
	BaseModel

	Name           string          `gorm:"type: varchar(50); not null" json:"name"`
	PropertyValues []PropertyValue `json:"property_values"`
}

func (PropertyName) TableName() string {
	return "property_name"
}
