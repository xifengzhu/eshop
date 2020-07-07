package models

type Region struct {
	BaseModel

	Name     string `gorm:"type: varchar(64); not null" json:"name"`
	RegionID string `gorm:"type: varchar(12); not null" json:"region_id"`
	CityID   string `gorm:"type: varchar(12); not null" json:"city_id"`
}
