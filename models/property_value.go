package models

type PropertyValue struct {
	BaseModel

	Value          string `gorm:"type: varchar(50); not null" json:"value"`
	PropertyNameID int    `gorm:"type: int; not null" json:"property_name_id"`
	Destroy        bool   `sql:"-" json:"_destroy,omitempty"`
}

func (PropertyValue) TableName() string {
	return "property_value"
}
