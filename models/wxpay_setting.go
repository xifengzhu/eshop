package models

type WxpaySetting struct {
	BaseModel

	WxappId   string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	AppName   string `gorm:"type: varchar(255); not null" json:"app_name"`
	AppSecret string `gorm:"type: varchar(255); not null" json:"app_secret"`
	Mchid     string `gorm:"type: varchar(255); not null" json:"mchid"`
	Apikey    string `gorm:"type: varchar(255); not null" json:"apikey"`
	NotifyUrl string `gorm:"type: varchar(255); " json:"notify_url"`
}

func (setting *WxpaySetting) CreateOrUpdate() (err error) {
	var existSetting WxpaySetting
	err = db.First(&existSetting).Error
	setting.ID = existSetting.ID
	err = db.Save(&setting).Error
	return
}

func (setting *WxpaySetting) Current() (err error) {
	err = db.First(&setting).Error
	return
}
