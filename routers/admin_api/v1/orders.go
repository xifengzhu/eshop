package v1

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
	"time"
)

type DeliveryOrderParams struct {
	CartItemIDs  []int  `json:"item_ids" binding:"required,gt=0"`
	AddressID    int    `json:"address_id" binding:"required"`
	ExpressID    int    `json:"express_id" binding:"required"`
	BuyerMessage string `json:"buyer_message"`
}

type QueryOrderParams struct {
	utils.Pagination
	Status          []int      `json:"q[status_in]"`
	Order_no        string     `json:"q[order_no_cont]"`
	Created_at_gteq *time.Time `json:"q[created_at_gteq]"`
	Created_at_lteq *time.Time `json:"q[created_at_lteq]"`
}

// @Summary 订单发货
// @Produce  json
// @Tags 后台订单管理
// @Param id query DeliveryOrderParams true "order id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/orders/{id}/deliver [post]
// @Security ApiKeyAuth
func DeliveryOrder(c *gin.Context) {

}

// @Summary 获取订单列表
// @Produce  json
// @Tags 后台订单管理
// @Param params query QueryOrderParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/orders [get]
// @Security ApiKeyAuth
func GetOrders(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Order
	result := &[]models.Order{}

	models.SearchResourceQuery(&model, result, &pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 获取订单详情
// @Produce  json
// @Tags 后台订单管理
// @Param id path int true "order id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/orders/{id} [get]
// @Security ApiKeyAuth
func GetOrder(c *gin.Context) {
	var order models.Order
	id, _ := strconv.Atoi(c.Param("id"))
	order.ID = id

	err := models.FindResource(&order, Query{Preloads: []string{"OrderItems"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, order)
}
