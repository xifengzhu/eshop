package v1

import (
	"github.com/gin-gonic/gin"
	. "github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/helpers"
)

// @Summary dashboard
// @Produce  json
// @Tags 数据统计
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/dashboard [get]
func Dashboard(c *gin.Context) {
	// TODO:
	data := make(map[string]interface{})
	dataOverview := make(map[string]interface{})
	dataOverview["wait_seller_send_goods_count"] = Recent7DayNewOrderCount()
	dataOverview["after_sale_count"] = 12
	dataOverview["total_order_amount"] = TotalOrderAmount()
	dataOverview["week_new_order_amount"] = Total7DayOrderAmount()
	dataOverview["total_user_count"] = TotalUserCount()
	dataOverview["today_new_user_count"] = TodayNewUserCount()
	dataOverview["today_new_order_amount"] = TotalTodayOrderAmount()
	dataOverview["today_new_order_count"] = TodayNewOrderCount()

	orderTrend := make(map[string]interface{})
	orderTrend["order_num_trend"] = OrderNumTrend()
	orderTrend["paid_order_num_trend"] = PaidOrderNumTrend()
	orderTrend["order_actual_amount_trend"] = ActualOrderAmountTrend()
	orderTrend["order_paid_amount_trend"] = PaidOrderAmountTrend()

	dataOverview["recent_30_order_trend"] = orderTrend

	data["data_overview"] = dataOverview
	apiHelpers.ResponseSuccess(c, data)
}
