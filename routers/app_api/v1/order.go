package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
	"github.com/xifengzhu/eshop/routers/app_api/entities"
	"log"
	"strconv"
	"time"
)

type QueryOrderParams struct {
	utils.Pagination
	status string `json:"q[status]"`
}

type OrderParams struct {
	AddressID    int    `json:"address_id" binding:"required"`
	ExpressID    int    `json:"express_id" binding:"required"`
	BuyerMessage string `json:"buyer_message"`
}

// @Summary 获取订单列表
// @Produce  json
// @Tags 订单
// @Param params query QueryOrderParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders [get]
// @Security ApiKeyAuth
func GetOrders(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Order
	orders := &[]models.Order{}

	condition := c.QueryMap("q")
	condition["user_id"] = strconv.Itoa(user.ID)

	models.SearchResourceQuery(&model, orders, pagination, condition)

	orderEntities := transferOrderToEntity(*orders)

	response := apiHelpers.Collection{Pagination: pagination, List: orderEntities}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 获取订单详情
// @Produce  json
// @Tags 订单
// @Param id path int true "order id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/{id} [get]
// @Security ApiKeyAuth
func GetOrder(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var userOrders []models.Order
	parmMap := map[string]interface{}{"id": orderID, "user_id": user.ID}
	err := models.WhereResources(&userOrders, Query{Conditions: parmMap, Preloads: []string{"OrderItems"}})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}
	apiHelpers.ResponseSuccess(c, userOrders[0])
}

// @Summary 下单预检查(待实现)
// @Produce  json
// @Tags 订单
// @Param params body OrderParams true "订单参数"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/pre_check [post]
// @Security ApiKeyAuth
func PreOrder(c *gin.Context) {
	var order models.Order
	apiHelpers.ResponseSuccess(c, order)
}

// @Summary 创建订单
// @Produce  json
// @Tags 订单
// @Param params body OrderParams true "订单参数"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders [post]
// @Security ApiKeyAuth
func CreateOrder(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	var orderParams OrderParams
	if err := c.ShouldBindJSON(&orderParams); err != nil {
		log.Println("===orderParans bind===", err)
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	address, err := user.GetAddressByID(orderParams.AddressID)
	receiverProperties := address.DisplayString()
	log.Println("===receiverProperties===", receiverProperties)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var orderItems []models.OrderItem

	carItems, _ := user.GetCheckedShoppingCartItems()

	for _, item := range carItems {
		var goods models.Goods
		goods.ID = item.GoodsID
		err = models.FindResource(&goods, Query{})
		log.Println("===goods.Find:===", goods)
		if err != nil {
			apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, errors.New("商品不存在或被下架"))
			return
		}
		orderItem := models.OrderItem{GoodsName: goods.Name, GoodsPrice: goods.Price, LinePrice: goods.LinePrice, GoodsWeight: goods.Weight, GoodsAttr: goods.PropertiesText, TotalNum: item.Quantity, DeductStockType: 10, GoodsID: item.GoodsID}
		orderItems = append(orderItems, orderItem)
	}

	log.Println("===orderIten===", orderItems)
	order := models.Order{
		WxappId:            "001",
		UserID:             user.ID,
		AddressID:          orderParams.AddressID,
		ExpressID:          orderParams.ExpressID,
		ReceiverProperties: address.DisplayString(),
		BuyerMessage:       orderParams.BuyerMessage,
		OrderItems:         orderItems,
	}

	err = models.CreateResource(&order)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	apiHelpers.ResponseSuccess(c, entities.OrderIDEntity{ID: order.ID})
}

// @Summary 请求支付参数
// @Produce  json
// @Tags 订单
// @Param order_id query integer true "订单ID"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/payment_params [post]
// @Security ApiKeyAuth
func RequestPayment(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	orderID, _ := strconv.Atoi(c.Query("order_id"))
	order, err := user.GetOrder(orderID)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, errors.New("订单不存在或已过期"))
		return
	}
	paymentParams, err := order.RequestPayment()
	if err != nil {
		apiHelpers.ResponseError(c, e.WECHAT_PAY_ERROR, err)
		return
	}
	apiHelpers.ResponseSuccess(c, paymentParams)
}

// @Summary 删除订单
// @Produce  json
// @Tags 订单
// @Param id path integer true "订单id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/{id} [delete]
// @Security ApiKeyAuth
func DeleteOrder(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var order models.Order
	parmMap := map[string]interface{}{"id": orderID, "user_id": user.ID}
	err := models.FirstResource(&order, Query{Conditions: parmMap})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, errors.New("资源不存在"))
		return
	}

	models.DestroyResource(&order, Query{Callbacks: []func(){order.DestroyOrderItems}})

	models.DestroyResource(order, Query{})
	apiHelpers.ResponseOK(c)
}

// @Summary 订单支付结果通知
// @Produce  json
// @Tags 微信回调通知
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/pay_notify [post]
func PayNotify(c *gin.Context) {
	params, _ := c.GetRawData()
	type NotifyParams struct {
		outTradeNo    string    `xml:"out_trade_no"`
		timeEnd       time.Time `xml:"time_end"`
		totalFee      string    `xml:"total_fee"`
		openID        string    `xml:"openid"`
		transactionID string    `xml:"transaction_id"`
	}
	log.Println("pay notify params:", string(params))
}

// @Summary 订单退款结果通知
// @Produce  json
// @Tags 微信回调通知
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/orders/refund_notify [post]
func RefundNotify(c *gin.Context) {
	params, _ := c.GetRawData()
	log.Println("refund notify params:", params)
}

func transferOrderToEntity(orders []models.Order) (orderEntities []entities.OrderEntity) {
	for _, d_order := range orders {
		var orderEntity entities.OrderEntity
		copier.Copy(&orderEntity, &d_order)
		orderEntities = append(orderEntities, orderEntity)
	}
	return
}
