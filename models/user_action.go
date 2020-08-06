package models

import (
// "github.com/xifengzhu/eshop/helpers/utils"
)

type UserAction struct {
	BaseModel

	Action    string `gorm:"type: varchar(50); not null" json:"action"` // view, collect, buy, search
	UserID    int    `gorm:"type: int;" json:"user_id"`
	ProductID int    `gorm:"type: int;" json:"product_id"`
	Rating    int    `gorm:"type: int;" json:"rating"`
}

// TODO
type Similarity struct {
	i int     `gorm:"type: int; primary_key" json:"i"`
	j int     `gorm:"type: int; primary_key" json:"j"`
	p float64 `gorm:"type: int;" json:"p"`
}
