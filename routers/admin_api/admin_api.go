package admin_api

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/middleware/jwt"
	// "github.com/xifengzhu/eshop/middleware/role"
	"github.com/xifengzhu/eshop/routers/admin_api/v1"
)

func InitAdminAPI(r *gin.Engine) {

	admin_apiv1 := r.Group("/admin_api/v1")

	admin_apiv1.POST("/sessions/login", v1.Login)
	admin_apiv1.POST("/sessions/forget_password", v1.ForgetPassword)
	admin_apiv1.PUT("/sessions/reset_password", v1.ResetPassword)
	admin_apiv1.GET("/get_captcha", v1.GetCaptcha)
	admin_apiv1.Use(jwt.JWTAuth())
	// admin_apiv1.Use(role.AuthCheckRole())
	{
		admin_apiv1.GET("/orders", v1.GetOrders)
		admin_apiv1.GET("/orders/:id", v1.GetOrder)
		admin_apiv1.POST("/orders/:id/ship", v1.ShipOrder)
		admin_apiv1.POST("/orders/:id/pay", v1.PayOrder)
	}

	{
		admin_apiv1.POST("/admin_users", v1.AddAdminUser)
		admin_apiv1.GET("/admin_users", v1.GetAdminUsers)
		admin_apiv1.GET("/admin_users/:id", v1.GetAdminUser)
		admin_apiv1.PUT("/admin_users/:id", v1.UpdateAdminUser)
		admin_apiv1.DELETE("/admin_users/:id", v1.DeleteAdminUser)
		admin_apiv1.GET("/sessions/mine", v1.GetCurrentAdminUsers)
		admin_apiv1.POST("/cabin_rules", v1.AddPolicy)
		admin_apiv1.DELETE("/cabin_rules", v1.RemovePolicy)
	}

	{
		admin_apiv1.POST("/categories", v1.AddProductGroup)
		admin_apiv1.GET("/categories/:id", v1.GetCategory)
		admin_apiv1.DELETE("/categories/:id", v1.DeleteCategory)
		admin_apiv1.GET("/categories", v1.GetCategories)
		admin_apiv1.PUT("/categories/:id", v1.UpdateCategory)
	}

	{
		admin_apiv1.POST("/product_groups", v1.AddProductGroup)
		admin_apiv1.GET("/product_groups/:id", v1.GetProductGroup)
		admin_apiv1.DELETE("/product_groups/:id", v1.DeleteProductGroup)
		admin_apiv1.GET("/product_groups", v1.GetProductGroups)
		admin_apiv1.PUT("/product_groups/:id", v1.UpdateProductGroup)
	}

	{
		admin_apiv1.GET("/app_setting", v1.GetAppSetting)
		admin_apiv1.PUT("/app_setting", v1.UpdateAppSetting)
		admin_apiv1.POST("/app_setting/cert", v1.UpdateWechatCert)
	}

	{
		admin_apiv1.POST("/deliveries", v1.AddDelivery)
		admin_apiv1.GET("/deliveries/:id", v1.GetDelivery)
		admin_apiv1.DELETE("/deliveries/:id", v1.DeleteDelivery)
		admin_apiv1.GET("/deliveries", v1.GetDeliveries)
		admin_apiv1.PUT("/deliveries/:id", v1.UpdateDelivery)
	}

	{
		admin_apiv1.POST("/expresses", v1.AddExpress)
		admin_apiv1.GET("/expresses/:id", v1.GetExpress)
		admin_apiv1.DELETE("/expresses/:id", v1.DeleteExpress)
		admin_apiv1.GET("/expresses", v1.GetExpresses)
		admin_apiv1.PUT("/expresses/:id", v1.UpdateExpress)
	}

	{
		admin_apiv1.POST("/products", v1.AddProduct)
		admin_apiv1.GET("/products/:id", v1.GetProduct)
		admin_apiv1.DELETE("/products/:id", v1.DeleteProduct)
		admin_apiv1.GET("/products", v1.GetProductes)
		admin_apiv1.PUT("/products/:id", v1.UpdateProduct)
	}

	{
		admin_apiv1.POST("/property_names", v1.AddPropertyName)
		admin_apiv1.GET("/property_names/:id", v1.GetPropertyName)
		admin_apiv1.DELETE("/property_names/:id", v1.DeletePropertyName)
		admin_apiv1.GET("/property_names", v1.GetPropertyNames)
		admin_apiv1.PUT("/property_names/:id", v1.UpdatePropertyName)

		admin_apiv1.POST("/property_values", v1.AddPropertyValue)
		admin_apiv1.DELETE("/property_values/:id", v1.DeletePropertyValue)
	}

	{
		admin_apiv1.POST("/web_pages", v1.AddWebPage)
		admin_apiv1.GET("/web_pages/:id", v1.GetWebPage)
		admin_apiv1.DELETE("/web_pages/:id", v1.DeleteWebPage)
		admin_apiv1.GET("/web_pages", v1.GetWebPages)
		admin_apiv1.PUT("/web_pages/:id", v1.UpdateWebPage)
	}

	{
		admin_apiv1.POST("/wxapp_pages", v1.AddWxAppPage)
		admin_apiv1.GET("/page_group_links", v1.GetPageGroupLinks)
		admin_apiv1.GET("/wxapp_pages/:id", v1.GetWxAppPage)
		admin_apiv1.DELETE("/wxapp_pages/:id", v1.DeleteWxAppPage)
		admin_apiv1.GET("/wxapp_pages", v1.GetWxAppPages)
		admin_apiv1.PUT("/wxapp_pages/:id", v1.UpdateWxAppPage)
	}

	{
		admin_apiv1.GET("/users/:id", v1.GetUser)
		admin_apiv1.GET("/users", v1.GetUsers)
	}

	{
		admin_apiv1.GET("/dashboard", v1.Dashboard)
		admin_apiv1.GET("/qiniu_meta", v1.GetQiniuMeta)
	}

	{
		admin_apiv1.GET("/addresses", v1.GetUserAdddresses)
		admin_apiv1.GET("/provinces", v1.GetProvinces)
		admin_apiv1.GET("/cities", v1.GetCities)
		admin_apiv1.GET("/regions", v1.GetRegions)
	}
}
