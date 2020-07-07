package models

import (
	"fmt"
)

type OrderItem struct {
	BaseModel

	WxappId         string  `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	GoodsName       string  `gorm:"type: varchar(255); not null" json:"goods_name"`
	Cover           string  `gorm:"type: varchar(255); not null" json:"cover"`
	GoodsPrice      float32 `gorm:"type: decimal(10,2);" json:"goods_price"`
	LinePrice       float32 `gorm:"type: decimal(10,2);" json:"line_price"`
	GoodsWeight     float32 `gorm:"type: decimal(10,2);" json:"goods_weight"`
	GoodsAttr       string  `gorm:"type: varchar(255);" json:"goods_attr"`
	TotalNum        int     `gorm:"type: int;" json:"total_num"`
	TotalAmount     float32 `gorm:"type: decimal(10,2);" json:"total_amount"`
	DeductStockType int     `gorm:"type: tinyint;" json:"deduct_stock_type"`

	OrderId int `gorm:"type: int;" json:"order_id"`
	GoodsID int `gorm:"type: int;" json:"goods_id"`
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

// Callbacks
func (orderItem *OrderItem) BeforeSave() (err error) {
	fmt.Println("=======order item before save=============")
	orderItem.caculateTotalAmount() // 计算订单项总金额
	return nil
}

func (orderItem *OrderItem) caculateTotalAmount() {
	orderItem.TotalAmount = orderItem.GoodsPrice * float32(orderItem.TotalNum)
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
