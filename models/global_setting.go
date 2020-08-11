package models

type GlobalSetting struct {
	BaseModel
	DeductStockType      int     `gorm:"type: tinyint; not null" json:"deduct_stock_type"`
	FreeFreightAmount    float64 `gorm:"type: decimal(10,2);" json:"free_freight_amount"`
	FreeFreightRegion    float64 `sql:"type: json;" json:"region,omitempty"`
	DefaultShareImg      string  `gorm:"type: varchar(255); not null" json:"default_share_img"`
	DefaultShareSentence string  `gorm:"type: varchar(255); not null" json:"default_share_sentence"`
	HotSearchWords       string  `sql:"type: json;" json:"hot_search_words"`
}

func (setting *GlobalSetting) CreateOrUpdate() (err error) {
	var existSetting GlobalSetting
	db.First(&existSetting)
	setting.ID = existSetting.ID
	err = db.Save(&setting).Error
	return
}

func (setting *GlobalSetting) Current() (err error) {
	err = db.First(&setting).Error
	return
}
