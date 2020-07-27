package entities

import (
	"github.com/xifengzhu/eshop/models"
	"time"
)

type OrderDetailEntity struct {
	OrderEntity
	OuterPayId         string          `json:"outer_pay_id"`
	PayAt              *time.Time      `json:"pay_at"`
	ReceiverProperties string          `json:"receiver_properties"`
	BuyerMessage       string          `json:"buyer_message"`
	Logistic           models.Logistic `json:"logistic"`
	UserID             int             `json:"user_id"`
}

type OrderEntity struct {
	ID                int                `json:"id"`
	WxappId           string             `json:"wxapp_id"`
	OrderNo           string             `json:"order_no"`
	State             string             `json:"state"`
	ExpressFee        float64            `json:"express_fee"`
	PayAmount         float64            `json:"pay_amount"`
	ProductAmount     float64            `json:"product_amount"`
	LatestPaymentTime *time.Time         `json:"latest_payment_time"`
	OrderItems        []models.OrderItem `json:"order_items"`
	Coupon            models.Coupon      `json:"coupon"`
	Coupons           []models.Coupon    `json:"coupons,omitempty"`
	Address           models.Address     `json:"address"`
}

type OrderIDEntity struct {
	ID int `json:"order_id"`
}
