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

	PropertiesText string `gorm:"-" json:"properties_text"`
	PIDs           []int  `gorm:"-" json:"pids"`
	VIDs           []int  `gorm:"-" json:"vids"`
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

func (goods *Goods) SetPropertiesText() {
	var text string
	properties := strings.Split(goods.Properties, ";")
	for _, property := range properties {
		values := strings.Split(property, ":")

		var pvalue PropertyValue
		pvID, _ := strconv.Atoi(values[1])
		pvalue.ID = pvID

		FindResource(&pvalue, Options{})

		text += pvalue.Name + " "
	}
	goods.PropertiesText = strings.Trim(text, "\t \n")
}

func (goods *Goods) SetPropertyIDs() {
	var pids []int
	var vids []int
	properties := strings.Split(goods.Properties, ";")
	for _, property := range properties {
		values := strings.Split(property, ":")
		pid, _ := strconv.Atoi(values[0])
		vid, _ := strconv.Atoi(values[1])
		pids = append(pids, pid)
		vids = append(vids, vid)
	}
	goods.PIDs = pids
	goods.VIDs = vids
}

func (goods *Goods) AfterFind() (err error) {
	goods.SetPropertiesText()
	return
}

func (goods Goods) DeliveryRule(expressID int, provinceID int) (rule DeliveryRule) {
	var product Product
	db.Where("id = ?", goods.ProductID).Find(&product)
	rule = product.DeliveryRule(expressID, provinceID)
	return
}
