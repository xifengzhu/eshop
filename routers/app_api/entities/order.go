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
	ExpressAmount     float32            `json:"express_amount"`
	PayAmount         float32            `json:"pay_amount"`
	TotalAmount       float32            `json:"total_amount"`
	LatestPaymentTime *time.Time         `json:"latest_payment_time"`
	OrderItems        []models.OrderItem `json:"order_items"`
}

type OrderIDEntity struct {
	ID int `json:"order_id"`
}
