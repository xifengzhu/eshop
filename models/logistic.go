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
	Trace          JSON   `gorm:"type: json; " json:"trace"`
	OrderID        int    `gorm:"type: int; " json:"order_id"`
	Order          *Order `json:"order"`
}

func (Logistic) TableName() string {
	return "logistic"
}

func (logistic *Logistic) AfterCreate(tx *gorm.DB) (err error) {
	var order Order
	order.ID = logistic.OrderID
	err = FindResource(&order, Options{})
	if err == nil {
		err = order.Ship(tx)
		if err != nil {
			log.Println("ship order err", err)
			return
		}
	}
	return
}
