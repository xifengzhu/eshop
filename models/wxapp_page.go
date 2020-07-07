package models

type WxappPage struct {
	BaseModel

	WxappId string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name    string `gorm:"type: varchar(50); not null" json:"name"`
	Key     string `gorm:"type: varchar(50); not null; unique" json:"key"`
	// page_type: 1 为默认页面不可删除，2 为自定义页面可删除
	PageType             string `gorm:"type: tinyint; not null" json:"page_type"`
	PageData             JSON   `gorm:"type: json;" json:"page_data"`
	ShareSentence        string `gorm:"type: text; " json:"share_sentence"`
	ShareCover           string `gorm:"type: text; " json:"share_cover"`
	ShareBackgroundCover string `gorm:"type: text; " json:"share_background_cover"`
}
