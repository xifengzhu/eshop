package app_api

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/middleware/jwt"
	"github.com/xifengzhu/eshop/routers/app_api/v1"
)

func InitAppAPI(r *gin.Engine) {
	apiv1 := r.Group("/app_api/v1")
	// get address data
	{
		apiv1.GET("/provinces", v1.GetProvinces)
		apiv1.GET("/cities", v1.GetCities)
		apiv1.GET("/regions", v1.GetRegions)
	}

	{
		apiv1.GET("/web_page", v1.GetWxappPage)
		apiv1.GET("/wxapp_page", v1.GetWxappPage)
	}

	// 回调通知接口
	{
		// 微信支付回调通知
		apiv1.POST("/orders/pay_notify", v1.PayNotify)
		apiv1.POST("/orders/refund_notify", v1.RefundNotify)
	}

	apiv1.GET("/user/auth", v1.AuthWithWechat)
	apiv1.GET("/user/fake_token", v1.GetToken)
	apiv1.POST("/user/verify", v1.VerifyToken)

	// product router
	{
		apiv1.GET("/products", v1.GetProducts)
		apiv1.GET("/batch_products", v1.BatchProducts)
		apiv1.GET("/products/:id", v1.GetProduct)
	}

	apiv1.Use(jwt.JWTAuth())
	// tag router
	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.GET("/tags/:id", v1.GetTag)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags", v1.DeleteTag)
	}
	// user router
	{
		apiv1.GET("/users/mine", v1.GetUser)
		apiv1.PUT("/users/update", v1.EditUser)
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
		apiv1.GET("/shopping_cart", v1.GetCartItems)
		apiv1.POST("/shopping_cart/add", v1.AddCartItem)
		apiv1.PUT("/shopping_cart/check", v1.CheckCartItem)
		apiv1.PUT("/shopping_cart/uncheck", v1.UnCheckCartItem)
		apiv1.DELETE("/shopping_cart/remove", v1.DeleteCartItem)
		apiv1.PUT("/shopping_cart/qty", v1.UpdateCartItemQty)
	}
	{
		apiv1.GET("/orders", v1.GetOrders)
		apiv1.POST("/orders", v1.CreateOrder)
		apiv1.GET("/orders/:id", v1.GetOrder)
		apiv1.POST("/orders/pre_check", v1.PreOrder)
		apiv1.POST("/orders/payment_params", v1.RequestPayment)
		apiv1.DELETE("/orders/:id", v1.DeleteOrder)
	}
}
