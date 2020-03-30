package models

import (
	"github.com/jinzhu/gorm"
	"log"
)

type Logistic struct {
	BaseModel

	ExpressCompany string `gorm:"type: varchar(50); not null" json:"express_company"`
	ExpressCode    string `gorm:"type: varchar(50); " json:"express_code"`
	ExpressNo      string `gorm:"type: varchar(255); not null" json:"express_no"`
	Trace          string `gorm:"type: text; " json:"trace"`
	OrderID        int    `gorm:"type: int; default 1" json:"order_id"`
}

func (Logistic) TableName() string {
	return "logistic"
}

func (logistic *Logistic) Order() (order Order, err error) {
	order.ID = logistic.OrderID
	err = FindResource(&order, Options{})
	return
}

func (logistic *Logistic) AfterCreate(tx *gorm.DB) (err error) {
	order, err := logistic.Order()
	err = order.Ship(tx)
	if err != nil {
		log.Println("ship order err", err)
		return
	}
	return
}
