package models

// import (
// 	"log"
// )

type OrderItem struct {
	BaseModel

	GoodsName        string        `gorm:"type: varchar(255); not null" json:"goods_name"`
	Cover            string        `gorm:"type: varchar(255); not null" json:"cover"`
	GoodsPrice       float64       `gorm:"type: decimal(10,2);" json:"goods_price"`
	LinePrice        float64       `gorm:"type: decimal(10,2);" json:"line_price"`
	GoodsWeight      float64       `gorm:"type: decimal(10,2);" json:"goods_weight"`
	GoodsAttr        string        `gorm:"type: varchar(255);" json:"goods_attr"`
	TotalNum         int           `gorm:"type: int;" json:"total_num"`
	ProductAmount    float64       `gorm:"type: decimal(10,2);" json:"product_amount"`
	TotalAmount      float64       `gorm:"type: decimal(10,2);" json:"total_amount"`
	AdjustmentAmount float64       `gorm:"type: decimal(10,2);" json:"adjustment_amount"`
	DeductStockType  int           `gorm:"type: tinyint;" json:"deduct_stock_type"`
	OrderId          int           `gorm:"type: int;" json:"order_id"`
	GoodsID          int           `gorm:"type: int;" json:"goods_id"`
	Adjustments      []*Adjustment `gorm:"polymorphic:Target;" json:"adjustments,omitempty"`
}

func (OrderItem) TableName() string {
	return "order_item"
}

func (orderItem OrderItem) GetGoods() (Goods, error) {
	var goods Goods
	goods.ID = orderItem.GoodsID
	err := Find(&goods, Options{})
	return goods, err
}

func (orderItem OrderItem) GetProduct() (product Product, err error) {
	goods, err := orderItem.GetGoods()
	if err != nil {
		return product, err
	}
	product.ID = goods.ProductID
	err = Find(&product, Options{})
	return
}

func (orderItem *OrderItem) CaculateProductAmount() {
	orderItem.ProductAmount = orderItem.GoodsPrice * float64(orderItem.TotalNum)
}

// 减库存
func (orderItem *OrderItem) ReductStock() (err error) {
	goods, _ := orderItem.GetGoods()
	err = goods.ReduceStock(orderItem.TotalNum)
	return
}

// 恢复库存
func (orderItem *OrderItem) RestoreStock() {
	goods, _ := orderItem.GetGoods()
	goods.IncreseStock(orderItem.TotalNum)
}
