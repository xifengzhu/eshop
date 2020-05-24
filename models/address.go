package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type Address struct {
	BaseModel
	UserID    int    `gorm:"type:int; not null" json:"user_id"`
	User      User   `gorm:"-" json:"-"`
	WxappId   string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Region    string `gorm:"type: varchar(100); not null" json:"region"`
	Province  string `gorm:"type: varchar(100); not null" json:"province"`
	City      string `gorm:"type: varchar(100); not null" json:"city"`
	Detail    string `gorm:"type: varchar(100); not null" json:"detail"`
	isDefault bool   `gorm:"type: boolean; DEFAULT:false" json:"is_default"`
	Phone     string `gorm:"type: varchar(20); not null" json:"phone"`
	Receiver  string `gorm:"type: varchar(50); not null" json:"receiver"`
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
	Receiver string `json:"receiver"`
}

func (address *Address) Exist() bool {
	if err := db.First(&address).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func (address Address) ToJSON() []byte {
	addressDisplay := AddressDisplay{Region: address.Region, Province: address.Province, City: address.City, Detail: address.Detail, Phone: address.Phone, Receiver: address.Receiver}
	addr, _ := json.Marshal(addressDisplay)
	return addr
}

func (address Address) DisplayString() string {
	return address.Province + address.City + address.Region + address.Detail + " " + address.Phone
}
