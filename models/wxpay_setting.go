package models

import (
	"encoding/base64"
	"log"
)

type WxpaySetting struct {
	BaseModel

	WxappId       string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	AppName       string `gorm:"type: varchar(255); not null" json:"app_name"`
	AppSecret     string `gorm:"type: varchar(255); not null" json:"app_secret"`
	ServicePhone  string `gorm:"type: varchar(50); " json:"service_phone"`
	Mchid         string `gorm:"type: varchar(255); not null" json:"mchid"`
	Apikey        string `gorm:"type: varchar(255); not null" json:"apikey"`
	NotifyUrl     string `gorm:"type: varchar(255); " json:"notify_url"`
	ApiClientCert string `gorm:"type: varchar(255);" json:"api_client_cert"`
}

func (WxpaySetting) TableName() string {
	return "wxpay_setting"
}

func (setting *WxpaySetting) CreateOrUpdate() (err error) {
	var existSetting WxpaySetting
	db.First(&existSetting)
	setting.ID = existSetting.ID
	err = db.Save(&setting).Error
	return
}

func (setting *WxpaySetting) Current() (err error) {
	err = db.First(&setting).Error
	return
}

func (setting *WxpaySetting) ApiClientCertData() (data string) {
	setting.Current()
	if setting.ID == 0 || setting.ApiClientCert == "" {
		return
	} else {
		base64Str := setting.ApiClientCert
		byteResult, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			log.Println("base64 encoding err:", err)
			return
		}
		data = string(byteResult[:])
		return
	}
}
