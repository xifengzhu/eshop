package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type Address struct {
	BaseModel
	UserID     int    `gorm:"type:int; not null" json:"user_id"`
	User       User   `gorm:"-" json:"-"`
	WxappId    string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	RegionID   int    `gorm:"type: int; not null" json:"region_id"`
	ProvinceID int    `gorm:"type: int; not null" json:"province_id"`
	CityID     int    `gorm:"type: int; not null" json:"city_id"`
	Detail     string `gorm:"type: varchar(100); not null" json:"detail"`
	isDefault  bool   `gorm:"type: boolean; DEFAULT:false" json:"is_default"`
	Phone      string `gorm:"type: varchar(20); not null" json:"phone"`
	Receiver   string `gorm:"type: varchar(50); not null" json:"receiver"`

	Region   *Region   `json:"region,omitempty"`
	Province *Province `json:"province,omitempty"`
	City     *City     `json:"city,omitempty"`
}

func (Address) TableName() string {
	return "address"
}

type AddressDisplay struct {
	Region   string `json:"region"`
	Province string `json:"province"`
	City     string `json:"city"`
	Detail   string `json:"detail"`
	Phone    string `json:"phone"`
}

func (address Address) Exist() bool {
	var result Address
	if err := db.First(&result, address.ID).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func (address Address) DisplayString() string {
	addressDisplay := AddressDisplay{Region: address.Region.Name, Province: address.Province.Name, City: address.City.Name, Detail: address.Detail, Phone: address.Phone}
	addr, _ := json.Marshal(addressDisplay)
	return string(addr)
}

func (address AddressDisplay) DisplayString() string {
	return address.Province + address.City + address.Region + address.Detail + " " + address.Phone
}
