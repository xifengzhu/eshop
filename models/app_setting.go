package models

type AppSetting struct {
	BaseModel

	WxappId      string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	AppName      string `gorm:"type: varchar(255); not null" json:"app_name"`
	AppSecret    string `gorm:"type: varchar(255); not null" json:"app_secret"`
	ServicePhone string `gorm:"type: varchar(50); not null" json:"service_phone"`
	Mchid        string `gorm:"type: varchar(255); not null" json:"mchid"`
	Apikey       string `gorm:"type: varchar(255); not null" json:"apikey"`
	NotifyUrl    string `gorm:"type: varchar(255); not null" json:"notify_url"`
}

func (AppSetting) TableName() string {
	return "app_setting"
}

func (setting *AppSetting) CreateOrUpdate() (err error) {
	var existSetting AppSetting
	db.First(&existSetting)
	setting.ID = existSetting.ID
	err = db.Save(&setting).Error
	return
}

func (setting *AppSetting) Current() (err error) {
	err = db.First(&setting).Error
	return
}
