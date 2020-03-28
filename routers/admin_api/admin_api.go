package admin_api

import (
	"github.com/gin-gonic/gin"
	// "github.com/xifengzhu/eshop/middleware/jwt"
	// "github.com/xifengzhu/eshop/middleware/role"
	"github.com/xifengzhu/eshop/routers/admin_api/v1"
)

func InitAdminAPI(r *gin.Engine) {

	admin_apiv1 := r.Group("/admin_api/v1")

	admin_apiv1.POST("/login", v1.Login)
	// admin_apiv1.Use(jwt.JWTAuth())
	// admin_apiv1.Use(role.AuthCheckRole())
	{
		admin_apiv1.POST("/orders", v1.DeliveryOrder)
		admin_apiv1.GET("/orders", v1.GetOrders)
		admin_apiv1.GET("/orders/:id", v1.GetOrder)
	}

	{
		admin_apiv1.POST("/admin_users", v1.AddAdminUser)
		admin_apiv1.GET("/admin_users", v1.GetAdminUsers)
		admin_apiv1.POST("/cabin_rules", v1.AddPolicy)
		admin_apiv1.DELETE("/cabin_rules", v1.RemovePolicy)
	}

	{
		admin_apiv1.POST("/categories", v1.AddCategory)
		admin_apiv1.GET("/categories/:id", v1.GetCategory)
		admin_apiv1.DELETE("/categories/:id", v1.DeleteCategory)
		admin_apiv1.GET("/categories", v1.GetCategories)
		admin_apiv1.PUT("/categories/:id", v1.UpdateCategory)
	}

	{
		admin_apiv1.GET("/app_setting", v1.GetAppSetting)
		admin_apiv1.PUT("/app_setting", v1.UpdateAppSetting)
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
}
