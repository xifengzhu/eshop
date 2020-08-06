package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type QueryOrderParams struct {
	utils.Pagination
	state string `json:"q[state]"`
}

type OrderParams struct {
	AddressID    int    `json:"address_id"`
	ExpressID    int    `json:"express_id"`
	CouponID     int    `json:"coupon_id"`
	BuyerMessage string `json:"buyer_message"`
	IsPreview    *bool  `json:"is_preview" validate:"required"`
}

type OrderIDParams struct {
	OrderID int `json:"order_id" validate:"required"`
}
