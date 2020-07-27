package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xifengzhu/eshop/helpers/utils"
	// "log"
	"time"
)

type Product struct {
	BaseModel

	WxappId        string          `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name           string          `gorm:"type: varchar(50); not null" json:"name"`
	SerialNo       string          `gorm:"type: varchar(15); not null; unique_index" json:"serial_no"`
	Cover          string          `gorm:"type: varchar(250);" json:"cover"`
	MainPictures   JSON            `gorm:"type: json;" json:"main_pictures"`
	ShareCover     string          `gorm:"type: varchar(250);" json:"share_cover"`
	ShareDesc      string          `gorm:"type: varchar(250);" json:"share_desc"`
	Content        string          `gorm:"type: text;" json:"content"`
	SalesInitial   int             `gorm:"type: int; default 1" json:"sales_initial"`
	SalesActual    int             `gorm:"type: int; default 1" json:"sales_actual"`
	Position       int             `gorm:"type: int; " json:"position"`
	Price          float64         `gorm:"type: decimal(10,2); " json:"price"`
	IsOnline       bool            `gorm:"type: boolean; default true" json:"is_online"`
	DeletedAt      *time.Time      `gorm:"type: datetime; " json:"deleted_at"`
	DeliveryID     int             `gorm:"type: int; " json:"delivery_id"`
	Goodses        []Goods         `json:"goodses,omitemptys"`
	Category       *Category       `json:"category,omitempty"`
	Delivery       *Delivery       `json:"delivery,omitempty"`
	Categories     []*Category     `gorm:"many2many:product_categories;" json:"categories"`
	Specifications []Specification `sql:"-" json:"specifications,omitempty"`
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
	var pids []int
	var propertyNames []PropertyName

	for index, _ := range product.Goodses {
		var gvids []int
		if index == 0 {
			pids, gvids = product.Goodses[index].PVIDs()
		} else {
			_, gvids = product.Goodses[index].PVIDs()
		}
		vids = append(vids, gvids...)
	}

	db.Model(&PropertyName{}).Preload("PropertyValues").Where(pids).Find(&propertyNames)

	for _, propertyName := range propertyNames {
		specification := Specification{Pid: propertyName.ID, Name: propertyName.Name}
		for _, propertyValue := range propertyName.PropertyValues {
			found := utils.ContainsInt(vids, int(propertyValue.ID))
			if found {
				pvalue := PValue{Vid: propertyValue.ID, Name: propertyValue.Value}
				specification.PValues = append(specification.PValues, pvalue)
			}
		}
		specifications = append(specifications, specification)
	}
	product.Specifications = specifications
}

func (product *Product) NestUpdate() (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		var goodses []Goods
		var deleteIDs []int
		for _, goods := range product.Goodses {
			if goods.Destroy == true {
				deleteIDs = append(deleteIDs, goods.ID)
			} else {
				goodses = append(goodses, goods)
			}
		}
		if len(deleteIDs) > 0 {
			err = tx.Where("id = ?", deleteIDs).Delete(Goods{}).Error
		}
		product.Goodses = goodses
		err = tx.Save(&product).Error
		return err
	})
}

func (p *Product) RemoveGoodses() {
	db.Where("product_id = ?", p.ID).Delete(Goods{})
}

func (product *Product) UpdateCatgories(categories []Category) {
	db.Model(&product).Association("Categories").Replace(categories)
}

func (product *Product) BeforeCreate() (err error) {
	for index, _ := range product.Goodses {
		product.Goodses[index].Name = product.Name
	}
	fmt.Println("======products BeforeCreate=====", product.Goodses)
	return
}

func (product *Product) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&Goods{}).Where("product_id = ?", product.ID).Update("name", product.Name)
	fmt.Println("======products AfterUpdate=====")
	return
}
