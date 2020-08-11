package models

type CarItem struct {
	BaseModel

	Quantity int    `gorm:"type: int; not null" json:"quantity"`
	Checked  bool   `gorm:"type: boolean; default false" json:"checked"`
	UserID   int    `gorm:"type: int; " json:"user_id"`
	GoodsID  int    `gorm:"type: int; not null" json:"goods_id"`
	Goods    *Goods `gorm:"association_autoupdate:false" json:"goods,omitempty"`
}
