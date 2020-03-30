package models

import (
	"fmt"
	"github.com/xifengzhu/eshop/helpers/utils"
	"log"
	"time"
)

type Product struct {
	BaseModel

	WxappId         string     `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name            string     `gorm:"type: varchar(50); not null" json:"name"`
	Content         string     `gorm:"type: text;" json:"content"`
	DeductStockType int        `gorm:"type: tinyint; not null" json:"deduct_stock_type"`
	SalesInitial    int        `gorm:"type: int; default 1" json:"sales_initial"`
	SalesActual     int        `gorm:"type: int; default 1" json:"sales_actual"`
	Position        int        `gorm:"type: int; " json:"position"`
	Price           float32    `gorm:"type: decimal(10,2); " json:"price"`
	IsOnline        bool       `gorm:"type: boolean; default true" json:"is_online"`
	DeletedAt       *time.Time `gorm:"type: datetime; " json:"deleted_at"`
	DeliveryID      int        `gorm:"type: int; " json:"delivery_id"`
	CategoryID      int        `gorm:"type: int; " json:"category_id"`
	Goodses         []Goods    `json:"goodses"`
	Category        Category   `json:"category"`
	Delivery        Delivery   `json:"delivery"`

	Specifications []Specification `gorm:"-" json:"specifications"`
}

func (Product) TableName() string {
	return "product"
}

type Specification struct {
	Pid     int      `json:"pid"`
	Name    string   `json:"name"`
	PValues []PValue `json:"p_values"`
}

type PValue struct {
	Vid  int    `json:"vid"`
	Name string `json:"name"`
}

func (product Product) DeliveryRule(expressID int, ProvinceID int) (rule DeliveryRule) {

	var rules []DeliveryRule
	db.Preload("Delivery").Where("delivery_id = ? AND express_id = ?", product.DeliveryID, expressID).Find(&rules)

	for _, tRule := range rules {
		if tRule.Hit(ProvinceID) {
			rule = tRule
			break
		}
	}
	fmt.Println("===hit rule====: ", rule)
	return rule
}

func (product *Product) SetSpecifications() {
	var specifications []Specification
	var vids []int
	var propertyNames []PropertyName

	for index, _ := range product.Goodses {
		product.Goodses[index].SetPropertyIDs()
		log.Println("goods.VIDs", product.Goodses[index].VIDs)
		vids = append(vids, product.Goodses[index].VIDs...)
	}

	pids := product.Goodses[0].PIDs

	db.Model(&PropertyName{}).Preload("PropertyValues").Where(pids).Find(&propertyNames)

	for _, propertyName := range propertyNames {
		specification := Specification{Pid: propertyName.ID, Name: propertyName.Name}
		for _, propertyValue := range propertyName.PropertyValues {
			found := utils.ContainsInt(vids, int(propertyValue.ID))
			if found {
				pvalue := PValue{Vid: propertyValue.ID, Name: propertyValue.Name}
				specification.PValues = append(specification.PValues, pvalue)
			}
		}
		specifications = append(specifications, specification)
	}
	product.Specifications = specifications
}
