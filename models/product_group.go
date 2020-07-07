package models

import (
	// "github.com/xifengzhu/eshop/helpers/utils"
	"encoding/json"
)

type ProductGroup struct {
	BaseModel

	WxappId    string    `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name       string    `gorm:"type: varchar(50); not null" json:"name"`
	Remark     string    `gorm:"type: varchar(100); not null" json:"remark"`
	ProductIDs JSON      `gorm:"type: json; " json:"product_ids"`
	Key        string    `gorm:"type: varchar(50); unique" json:"key"`
	Products   []Product `sql:"-" json:"products,omitempty"`
}

func (p *ProductGroup) GetProducts() (products []Product) {
	var productIDs []int
	json.Unmarshal([]byte(p.ProductIDs), &productIDs)
	db.Where(productIDs).Find(&products)
	return
}
