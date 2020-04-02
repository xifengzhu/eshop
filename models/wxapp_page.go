package models

type WxappPage struct {
	BaseModel

	WxappId string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name    string `gorm:"type: varchar(50); not null" json:"name"`
	// page_type: 1 为默认页面不可删除，2 为自定义页面可删除
	PageType string `gorm:"type: tinyint; not null" json:"page_type"`
	PageData string `gorm:"type: text; not null" json:"page_data"`
}

func (WxappPage) TableName() string {
	return "wxapp_page"
}
