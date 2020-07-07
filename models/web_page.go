package models

type WebPage struct {
	BaseModel

	WxappId string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Cover   string `gorm:"type: varchar(200); " json:"cover"`
	Title   string `gorm:"type: varchar(120); not null" json:"title"`
	Content string `gorm:"type: text; not null" json:"content"`
}
