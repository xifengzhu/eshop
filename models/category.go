package models

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type Category struct {
	BaseModel

	WxappId  string     `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	Name     string     `gorm:"type: varchar(50); not null" json:"name"`
	Position int        `gorm:"type: int; " json:"position"`
	ParentID int        `gorm:"type: int; DEFAULT: 0" json:"parent_id"`
	Image    string     `gorm:"type: varchar(100);" json:"image"`
	Children []Category `gorm:"foreignkey:parent_id" json:"children,omitempty"`
	Parent   *Category  `json:"parent,omitempty"`
	Products []*Product `gorm:"many2many:product_categories;" json:"products,omitempty"`
}

func (category Category) IsParent() bool {
	return category.ParentID == 0
}

// TODO: 支持分页
func (category Category) GetCategoryProducts(pagination *utils.Pagination) (products []Product) {
	if category.IsParent() {
		var categories []Category
		db.Model(&category).Association("Children").Find(&categories)
		db.Model(&categories).Related(&products, "Products")
	} else {
		db.Model(&category).Related(&products, "Products")
	}
	return
}

func (category Category) RemoveChildrenRefer() {
	db.Model(Category{}).Where("parent_id = ?", category.ID).UpdateColumn("parent_id", nil)
}
