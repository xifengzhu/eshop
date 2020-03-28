package role

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/models"
	"net/http"
)

//权限检查中间件
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		//根据上下文获取载荷claims 从claims获得role
		resource := c.MustGet("resource")
		role := resource.(models.AdminUser).Role
		//检查权限
		res, err := models.Enforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": -1,
				"msg":  "错误消息" + err.Error(),
			})
			return
		}
		if res {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "很抱歉您没有此权限",
			})
			return
		}
	}
}
