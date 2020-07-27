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

	"github.com/xifengzhu/eshop/helpers/export"
	"log"
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

	err = models.Find(&order, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	var ship ShipOrderParams
	if err = apiHelpers.ValidateParams(c, &ship); err != nil {
		return
	}

	logistic := models.Logistic{
		OrderID:        order.ID,
		ExpressCompany: ship.ExpressCompany,
		ExpressNo:      ship.ExpressNo,
	}

	err = models.Save(&logistic)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q"), Preloads: []string{"OrderItems", "User"}}, &result)

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

	err := models.Find(&order, Query{Preloads: []string{"OrderItems", "User", "Express"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
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

	err := models.Find(&order, Query{Preloads: []string{"Coupons"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	err = order.Pay()
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// reload
	models.Find(&order, Query{})
	apiHelpers.ResponseSuccess(c, order)
}

func ExportOrders(c *gin.Context) {
	var order models.Order

	filename, err := order.Export()

	if err != nil {
		log.Println("======export orders=======", err)
	}

	data := map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	}
	apiHelpers.ResponseSuccess(c, data)
}
