package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/helpers"
)

// @Summary 获得所有权限列表
// @Produce  json
// @Tags 角色权限管理
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/permissions [get]
// @Security ApiKeyAuth
func GetPermissions(c *gin.Context) {
	results := models.AllPermissions()
	ResponseSuccess(c, results)
}
