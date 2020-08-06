package v1

import (
	// "errors"
	"github.com/gin-gonic/gin"
	// "github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/export"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"log"
	"strconv"
)

// @Summary 订单发货
// @Produce  json
// @Tags 后台订单管理
// @Param id path integer true "order id"
// @Param params body params.ShipOrderParams true "发货参数"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/orders/{id}/ship [post]
// @Security ApiKeyAuth
func ShipOrder(c *gin.Context) {
	var order Order
	var err error
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err = Find(&order, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	var ship ShipOrderParams
	if err = ValidateParams(c, &ship, "json"); err != nil {
		return
	}

	logistic := Logistic{
		OrderID:        order.ID,
		ExpressCompany: ship.ExpressCompany,
		ExpressNo:      ship.ExpressNo,
	}

	err = Save(&logistic)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, logistic)
}

// @Summary 获取订单列表
// @Produce  json
// @Tags 后台订单管理
// @Param params query params.QueryOrderParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/orders [get]
// @Security ApiKeyAuth
func GetOrders(c *gin.Context) {
	pagination := SetDefaultPagination(c)

	var model Order
	result := &[]Order{}

	Search(&model, &SearchParams{
		Pagination: pagination,
		Conditions: c.QueryMap("q"),
		Preloads: []string{
			"OrderItems",
			"User",
		},
	}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 获取订单详情
// @Produce  json
// @Tags 后台订单管理
// @Param id path int true "order id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/orders/{id} [get]
// @Security ApiKeyAuth
func GetOrder(c *gin.Context) {
	var order Order
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err := Find(&order, Options{Preloads: []string{"OrderItems", "User", "Express"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, order)
}

// @Summary 后台支付订单
// @Produce  json
// @Tags 后台订单管理
// @Param id path int true "order id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/orders/{id}/pay [post]
// @Security ApiKeyAuth
func PayOrder(c *gin.Context) {
	var order Order
	order.ID, _ = strconv.Atoi(c.Param("id"))

	err := Find(&order, Options{Preloads: []string{"Coupons"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	err = order.Pay()
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	// reload
	Find(&order, Options{})
	ResponseSuccess(c, order)
}

func ExportOrders(c *gin.Context) {
	var order Order

	filename, err := order.Export()

	if err != nil {
		log.Println("======export orders=======", err)
	}

	data := map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	}
	ResponseSuccess(c, data)
}
