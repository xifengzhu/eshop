package models

type WebPage struct {
	BaseModel

	WxappId string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Title   string `gorm:"type: varchar(120); not null" json:"title"`
	Content string `gorm:"type: text; not null" json:"Content"`
}

func (WebPage) TableName() string {
	return "web_page"
}

func (webPage *WebPage) FindByTitle(title string) (err error) {
	err = db.Where("title = ?", title).First(&webPage).Error
	return
}
