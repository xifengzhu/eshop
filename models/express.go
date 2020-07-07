package models

import (
// "github.com/xifengzhu/eshop/helpers/utils"
)

type Express struct {
	BaseModel

	Name string `gorm:"type: varchar(50); not null" json:"name"`
	Code string `gorm:"type: varchar(10); not null" json:"code"`
}
