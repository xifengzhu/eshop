package models

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type Category struct {
	BaseModel

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

func (category Category) GetCategoryProducts(pagination *utils.Pagination) (products []Product) {
	offset := (pagination.Page - 1) * pagination.PerPage
	var categoryIDs []int
	if category.IsParent() {
		db.Model(&Category{}).Where("parent_id = ?", category.ID).Pluck("id", &categoryIDs)
	} else {
		categoryIDs = append(categoryIDs, category.ID)
	}
	joinQuery := db.Table("product").Joins("INNER JOIN product_categories ON product.id = product_categories.product_id AND product_categories.category_id IN (?)", categoryIDs)
	joinQuery.Select("count(distinct(id))").Count(&pagination.Total)
	joinQuery.Limit(pagination.PerPage).Offset(offset).Scan(&products)
	return
}

func (category Category) RemoveChildrenRefer() {
	db.Model(Category{}).Where("parent_id = ?", category.ID).UpdateColumn("parent_id", nil)
}
