package v1

import (
	"github.com/gin-gonic/gin"
	// "github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary dashboard
// @Produce  json
// @Tags 数据统计
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/dashboard [get]
func Dashboard(c *gin.Context) {
	data := make(map[string]interface{})
	dataOverview := make(map[string]interface{})
	dataOverview["wait_seller_send_goods_count"] = 12
	dataOverview["after_sale_count"] = 12
	dataOverview["total_order_amount"] = 23.90
	dataOverview["week_new_order_amount"] = 12.34
	dataOverview["total_user_count"] = 12
	dataOverview["today_new_user_count"] = 12
	dataOverview["today_new_order_count"] = 12
	data["data_overview"] = dataOverview
	apiHelpers.ResponseSuccess(c, data)
}
