package models

type WebPage struct {
	BaseModel

	Cover   string `gorm:"type: varchar(200); " json:"cover"`
	Title   string `gorm:"type: varchar(120); not null" json:"title"`
	Content string `gorm:"type: text; not null" json:"content"`
}
