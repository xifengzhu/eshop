package app_api

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/initializers/setting"
	. "github.com/xifengzhu/eshop/middlewares"
	"github.com/xifengzhu/eshop/routers/app_api/v1"
	"net/http"
)

func InitAppAPI(r *gin.Engine) {

	api_public := r.Group("/public")
	api_public.StaticFS("/qrcode", http.Dir(qrcodePath()))

	apiv1 := r.Group("/app_api/v1")
	apiv1.Use(IPFilter())
	// get address data
	{
		apiv1.GET("/provinces", v1.GetProvinces)
	}

	{
		apiv1.GET("/web_page", v1.GetWebPage)
		apiv1.GET("/wxapp_page", v1.GetWxappPage)
	}

	{
		apiv1.GET("/wechat/wxacode", v1.GetWxaCode)
	}

	// 回调通知接口
	{
		// 微信支付回调通知
		apiv1.POST("/orders/pay_notify", v1.PayNotify)
		apiv1.POST("/orders/refund_notify", v1.RefundNotify)
	}

	apiv1.POST("/user/auth", v1.AuthWithWechat)
	apiv1.GET("/user/fake_token", v1.GetToken)
	apiv1.POST("/user/verify", v1.VerifyToken)

	// product router
	{
		apiv1.GET("/products", v1.GetProducts)
		apiv1.GET("/recommend_products", v1.GetRecommendProducts)
		apiv1.GET("/batch_products", v1.BatchProducts)
		apiv1.GET("/products/:id", v1.GetProduct)
	}

	// categories
	{
		apiv1.GET("/categories", v1.GetCategories)
		apiv1.GET("/categories/:id/products", v1.GetCategoryProducts)
	}

	apiv1.Use(JWTAuth())
	// user router
	{
		apiv1.GET("/users/mine", v1.GetUser)
		apiv1.PUT("/users/mine", v1.EditUser)
	}
	// address router
	{
		apiv1.GET("/addresses", v1.GetAddresses)
		apiv1.GET("/addresses/:id", v1.GetAddress)
		apiv1.POST("/addresses", v1.AddAddress)
		apiv1.PUT("/addresses/:id", v1.EditAddress)
		apiv1.DELETE("/addresses/:id", v1.DeleteAddress)
	}
	// shopping cart
	{
		apiv1.GET("/shopping_cart/my", v1.GetCartItems)
		apiv1.POST("/shopping_cart/add", v1.AddCartItem)
		apiv1.PUT("/shopping_cart/check", v1.CheckCartItem)
		apiv1.PUT("/shopping_cart/uncheck", v1.UnCheckCartItem)
		apiv1.DELETE("/shopping_cart/delete", v1.DeleteCartItem)
		apiv1.PUT("/shopping_cart/qty", v1.UpdateCartItemQty)
	}
	// user orders
	{
		apiv1.GET("/orders", v1.GetOrders)
		apiv1.POST("/orders", v1.CreateOrder)
		apiv1.GET("/orders/:id", v1.GetOrder)
		apiv1.POST("/orders/request_payment", v1.RequestPayment)
		apiv1.POST("/orders/close", v1.CloseOrder)
		apiv1.DELETE("/orders/:id", v1.DeleteOrder)
	}

	{
		apiv1.POST("/coupons/receive", v1.CaptchCoupon)
		apiv1.GET("/coupons", v1.GetCoupons)
	}
}

func qrcodePath() string {
	return setting.RuntimeRootPath + "/public/qrcode/"
}
