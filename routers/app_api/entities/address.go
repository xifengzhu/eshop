package entities

import (
	"time"
)

type AddressPresent struct {
	ID         int   `json:"id"`
	UserID     int   `json:"user_id"`
	WxappId    string `json:"wxapp_id"`
	RegionID   int   `json:"region_id"`
	ProvinceID int   `json:"province_id"`
	CityID     int   `json:"city_id"`
	Detail     string `json:"detail"`
	isDefault  bool   `json:"is_default"`
	Phone      string `json:"phone"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	RegionName   string `json:"region"`
	ProvinceName string `json:"province"`
	CityName     string `json:"city"`
}
