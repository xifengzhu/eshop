package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/app_api/helpers"
	. "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/app_api/present"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"log"
	"strconv"
	"time"
)

// @Summary 获取订单列表
// @Produce  json
// @Tags 订单
// @Param params query params.QueryOrderParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders [get]
// @Security ApiKeyAuth
func GetOrders(c *gin.Context) {
	user := CurrentUser(c)

	pagination := SetDefaultPagination(c)

	var model Order
	orders := &[]Order{}

	condition := c.QueryMap("q")
	condition["user_id"] = strconv.Itoa(user.ID)

	Search(&model, &SearchParams{Pagination: pagination, Conditions: condition, Preloads: []string{"OrderItems"}}, orders)

	orderEntities := transferOrdersToEntity(*orders)

	response := Collection{Pagination: pagination, List: orderEntities}

	ResponseSuccess(c, response)
}

// @Summary 获取订单详情
// @Produce  json
// @Tags 订单
// @Param id path int true "order id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders/{id} [get]
// @Security ApiKeyAuth
func GetOrder(c *gin.Context) {
	user := CurrentUser(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	order, err := user.GetOrder(orderID)
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	orderEntity := transferOrderToEntity(order)

	ResponseSuccess(c, orderEntity)
}

// @Summary 创建订单
// @Produce  json
// @Tags 订单
// @Param params body params.OrderParams true "订单参数"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders [post]
// @Security ApiKeyAuth
func CreateOrder(c *gin.Context) {
	user := CurrentUser(c)
	var err error
	var orderParams OrderParams

	if err = ValidateParams(c, &orderParams, "json"); err != nil {
		return
	}

	var orderItems []OrderItem

	carItems, _ := user.GetCheckedShoppingCartItems()

	for _, item := range carItems {
		var goods Goods
		goods.ID = item.GoodsID
		err = Find(&goods, Options{})
		if err != nil {
			ResponseError(c, e.ERROR_NOT_EXIST, "商品不存在或被下架")
			return
		}
		orderItem := OrderItem{
			GoodsName:       goods.Name,
			GoodsPrice:      goods.Price,
			LinePrice:       goods.LinePrice,
			GoodsWeight:     goods.Weight,
			GoodsAttr:       goods.PropertiesText,
			TotalNum:        item.Quantity,
			DeductStockType: 10,
			GoodsID:         item.GoodsID,
			Cover:           goods.Image,
		}
		orderItems = append(orderItems, orderItem)
	}

	order := Order{
		UserID:       user.ID,
		User:         &user,
		ExpressID:    orderParams.ExpressID,
		BuyerMessage: orderParams.BuyerMessage,
		OrderItems:   orderItems,
	}

	var address Address
	if orderParams.AddressID != 0 {
		address, _ = user.GetAddressByID(orderParams.AddressID)
		receiverProperties := address.DisplayString()
		order.AddressID = address.ID
		order.ReceiverProperties = receiverProperties
	}

	// 检查coupon
	var coupon Coupon
	if orderParams.CouponID != 0 {
		coupon.ID = orderParams.CouponID
		err = Find(&coupon, Options{})
		if err != nil || coupon.State != "actived" {
			ResponseError(c, e.ERROR_NOT_EXIST, "invalid coupon")
			return
		}
		order.CouponID = coupon.ID
	}

	if *orderParams.IsPreview {
		order.Address = &address
		order.Coupon = &coupon
		order.PreOrder()

		orderEntity := transferOrderToEntity(order)
		coupons := order.AvaliableCoupons()

		log.Println("======avaliable coupons:===", coupons)
		orderEntity.Coupons = coupons

		ResponseSuccess(c, orderEntity)
	} else {
		if order.AddressID == 0 {
			ResponseError(c, e.ERROR_NOT_EXIST, "地址不能为空")
			return
		}

		err = Create(&order)
		if err != nil {
			ResponseError(c, e.INVALID_PARAMS, err.Error())
			return
		}
		ResponseSuccess(c, OrderIDEntity{ID: order.ID})
	}
}

// @Summary 请求支付参数
// @Produce  json
// @Tags 订单
// @Param params body params.OrderIDParams true "订单ID"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders/request_payment [post]
// @Security ApiKeyAuth
func RequestPayment(c *gin.Context) {
	user := CurrentUser(c)
	var params OrderIDParams
	if err := ValidateParams(c, &params, "json"); err != nil {
		return
	}

	order, err := user.GetOrder(params.OrderID)
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, "订单不存在或已过期")
		return
	}
	paymentParams, err := order.RequestPayment(c.ClientIP())
	if err != nil {
		ResponseError(c, e.WECHAT_PAY_ERROR, err.Error())
		return
	}
	ResponseSuccess(c, paymentParams)
}

// @Summary 删除订单
// @Produce  json
// @Tags 订单
// @Param id path integer true "订单id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders/{id} [delete]
// @Security ApiKeyAuth
func DeleteOrder(c *gin.Context) {
	user := CurrentUser(c)
	orderID, _ := strconv.Atoi(c.Param("id"))

	var order Order
	parmMap := map[string]interface{}{"id": orderID, "user_id": user.ID}
	err := Find(&order, Options{Conditions: parmMap})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, "资源不存在")
		return
	}

	DestroyWithCallbacks(&order, Options{Callbacks: []func(){order.DestroyOrderItems}})

	Destroy(order)
	ResponseOK(c)
}

// @Summary 取消订单
// @Produce  json
// @Tags 订单
// @Param params body params.OrderIDParams true "订单id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders/close [post]
// @Security ApiKeyAuth
func CloseOrder(c *gin.Context) {
	user := CurrentUser(c)

	var params OrderIDParams
	if err := ValidateParams(c, &params, "json"); err != nil {
		return
	}

	order, err := user.GetOrder(params.OrderID)

	err = order.Close()

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	ResponseOK(c)
}

// @Summary 订单支付结果通知
// @Produce  json
// @Tags 微信回调通知
// @Success 200 {object} helpers.Response
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
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/orders/refund_notify [post]
func RefundNotify(c *gin.Context) {
	params, _ := c.GetRawData()
	log.Println("refund notify params:", params)
}

func transferOrdersToEntity(orders []Order) (orderEntities []OrderEntity) {
	for _, d_order := range orders {
		var orderEntity OrderEntity
		copier.Copy(&orderEntity, &d_order)
		orderEntities = append(orderEntities, orderEntity)
	}
	return
}

func transferOrderToEntity(order Order) (orderEntity OrderDetailEntity) {
	copier.Copy(&orderEntity, &order)
	return
}
