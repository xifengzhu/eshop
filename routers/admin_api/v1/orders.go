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

type ShipOrderParams struct {
	ExpressNo string `json:"express_no" binding:"required`
}

type QueryOrderParams struct {
	utils.Pagination
	UserID          int        `json:"q[user_id_eq]"`
	State           []string   `json:"q[state_in]"`
	Order_no        string     `json:"q[order_no_cont]"`
	Created_at_gteq *time.Time `json:"q[created_at_gteq]"`
	Created_at_lteq *time.Time `json:"q[created_at_lteq]"`
}

// @Summary 订单发货
// @Produce  json
// @Tags 后台订单管理
// @Param id path integer true "order id"
// @Param params body ShipOrderParams true "发货参数"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/orders/{id}/ship [post]
// @Security ApiKeyAuth
func ShipOrder(c *gin.Context) {
	var order models.Order
	var err error
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err = models.FindResource(&order, Query{Preloads: []string{"Express"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	var ship ShipOrderParams
	if err = c.ShouldBindJSON(&ship); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	logistic := models.Logistic{
		OrderID:        order.ID,
		ExpressCompany: order.Express.Name,
		ExpressCode:    order.Express.Code,
		ExpressNo:      ship.ExpressNo,
	}

	err = models.SaveResource(&logistic)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, logistic)
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

	models.SearchResourceWithPreloadQuery(&model, result, pagination, c.QueryMap("q"), []string{"OrderItems", "User"})

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
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err := models.FindResource(&order, Query{Preloads: []string{"OrderItems", "User", "Express"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, order)
}

// @Summary 后台支付订单
// @Produce  json
// @Tags 后台订单管理
// @Param id path int true "order id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/orders/{id}/pay [post]
// @Security ApiKeyAuth
func PayOrder(c *gin.Context) {
	var order models.Order
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err := models.FindResource(&order, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}
	err = order.Pay()
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	// reload
	models.FindResource(&order, Query{})
	apiHelpers.ResponseSuccess(c, order)
}
