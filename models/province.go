package models

type Province struct {
	BaseModel

	Name       string `gorm:"type: varchar(64); not null" json:"name"`
	ProvinceID string `gorm:"type: varchar(25); not null" json:"province_id"`
}
