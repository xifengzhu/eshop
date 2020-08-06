package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 获得所有角色
// @Produce  json
// @Tags 角色权限管理
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/roles [get]
// @Security ApiKeyAuth
func GetRoles(c *gin.Context) {
	var roles []Role
	All(&roles, Options{})
	ResponseSuccess(c, roles)
}

// @Summary 添加权限给角色
// @Produce  json
// @Tags 角色权限管理
// @Param id path int true "role id"
// @Param params body params.AddPermisisonsParams true "权限"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/roles/{id}/permissions [post]
// @Security ApiKeyAuth
func AddPermissionToRole(c *gin.Context) {
	var role Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = id
	err := Find(&role, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	var params AddPermisisonsParams
	if err = ValidateParams(c, &params, "json"); err != nil {
		return
	}

	for _, permit := range params.Permissions {
		Enforcer.AddPermissionForUser(role.AuthKey(), permit)
	}
	ResponseOK(c)
}

// @Summary 角色详情
// @Produce  json
// @Tags 角色权限管理
// @Param id path int true "role id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/roles/{id} [get]
// @Security ApiKeyAuth
func GetRole(c *gin.Context) {
	var role Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = id
	err := Find(&role, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	permissions := role.GetPermissions()
	role.Permissions = permissions
	ResponseSuccess(c, role)
}
