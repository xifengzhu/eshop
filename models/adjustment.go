package models

type Adjustment struct {
	BaseModel

	OrderID    int     `gorm:"type: int"; json:"order_id"`
	Amount     float64 `gorm:"type: decimal(10,2);" json:"amount"`
	Label      string  `gorm:"type: varchar(100); json:"label"`
	SourceType string  `gorm:"type: varchar(20); json:"source_type"`
	SourceID   int     `gorm:"type: int"; json:"source_id"`
	TargetType string  `gorm:"type: varchar(20); json:"target_type"`
	TargetID   int     `gorm:"type: int"; json:"target_id"`
	Order      *Order  `json:"order,omitempty"`
}
