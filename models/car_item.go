package models

type CarItem struct {
	BaseModel

	WxappId  string `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Quantity int    `gorm:"type: int; not null" json:"quantity"`
	Checked  bool   `gorm:"type: boolean; default false" json:"checked"`
	UserID   int    `gorm:"type: int; " json:"user_id"`
	GoodsID  int    `gorm:"type: int; not null" json:"goods_id"`
	Goods    *Goods `gorm:"association_autoupdate:false" json:"goods,omitempty"`
}

func (CarItem) TableName() string {
	return "car_item"
}
