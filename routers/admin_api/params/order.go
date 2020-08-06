package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type ShipOrderParams struct {
	ExpressNo      string `json:"express_no" validate:"required`
	ExpressCompany string `json:"express_company" validate:"required`
}

type QueryOrderParams struct {
	utils.Pagination
	UserID          int        `json:"q[user_id_eq]"`
	State           []string   `json:"q[state_in]"`
	Order_no        string     `json:"q[order_no_cont]"`
	Created_at_gteq *time.Time `json:"q[created_at_gteq]"`
	Created_at_lteq *time.Time `json:"q[created_at_lteq]"`
}
