package models

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type Goods struct {
	BaseModel

	WxappId    string  `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name       string  `gorm:"type: varchar(120); not null" json:"name"`
	Properties string  `gorm:"type: varchar(255); not null" json:"properties"`
	Images     string  `gorm:"type: text; " json:"images"`
	SkuNo      string  `gorm:"type: varchar(50); not null" json:"sku_no"`
	StockNum   int     `gorm:"type: int; default 1" json:"stock_num"`
	Position   int     `gorm:"type: int; " json:"position"`
	Price      float32 `gorm:"type: decimal(10,2); " json:"price"`
	LinePrice  float32 `gorm:"type: decimal(10,2); " json:"line_price"`
	Weight     float32 `gorm:"type: double; " json:"weight"`
	ProductID  int     `gorm:"type: int; " json:"product_id"`
	Destroy    bool    `sql:"-" json:"_destroy,omitempty"`
}

func (Goods) TableName() string {
	return "goods"
}

func (goods *Goods) IsExist() bool {
	if err := db.Joins("LEFT JOIN product ON goods.product_id = product.id").Where("product.is_online IS TRUE").First(&goods).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func (goods *Goods) PropertiesText() string {
	var text string
	properties := strings.Split(goods.Properties, ";")
	for _, property := range properties {
		values := strings.Split(property, ":")

		var pvalue PropertyValue
		pvID, _ := strconv.Atoi(values[1])
		pvalue.ID = pvID

		FindResource(&pvalue, Options{})

		text += pvalue.Value + " "
	}
	return strings.Trim(text, "\t \n")
}

func (goods *Goods) PVIDs() (pids []int, vids []int) {
	properties := strings.Split(goods.Properties, ";")
	for _, property := range properties {
		values := strings.Split(property, ":")
		pid, _ := strconv.Atoi(values[0])
		vid, _ := strconv.Atoi(values[1])
		pids = append(pids, pid)
		vids = append(vids, vid)
	}
	return pids, vids
}

func (goods *Goods) AfterFind() (err error) {
	return
}

func (goods Goods) DeliveryRule(expressID int, provinceID int) (rule DeliveryRule) {
	var product Product
	db.Where("id = ?", goods.ProductID).Find(&product)
	rule = product.DeliveryRule(expressID, provinceID)
	return
}
