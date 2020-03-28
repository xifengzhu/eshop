package models

import (
	"time"
)

type Logistic struct {
	BaseModel

	ExpressCompany string     `gorm:"type: varchar(50); not null" json:"express_company"`
	ExpressCode    string     `gorm:"type: varchar(50); " json:"express_code"`
	ExpressNo      string     `gorm:"type: varchar(255); not null" json:"express_no"`
	DeliveryTime   *time.Time `gorm:"type: datetime; " json:"delivery_time"`
	Trace          string     `gorm:"type: text; " json:"trace"`
	OrderID        int       `gorm:"type: int; default 1" json:"order_id"`
}
