package models

type City struct {
	BaseModel

	Name       string `gorm:"type: varchar(64); not null" json:"name"`
	ProvinceID string `gorm:"type: varchar(12); not null" json:"province_id"`
	CityID     string `gorm:"type: varchar(12); not null" json:"city_id"`
}
