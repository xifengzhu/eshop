package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"log"
	"strconv"
)

type AddPermisisonsParams struct {
	Permissions []string `json:"permissions"`
}

// @Summary 获得所有角色
// @Produce  json
// @Tags 角色权限管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/roles [get]
// @Security ApiKeyAuth
func GetRoles(c *gin.Context) {
	var roles []models.Role
	models.All(&roles, Query{})
	apiHelpers.ResponseSuccess(c, roles)
}

// @Summary 添加权限给角色
// @Produce  json
// @Tags 角色权限管理
// @Param id path int true "role id"
// @Param params body AddPermisisonsParams true "权限"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/roles/{id}/permissions [post]
// @Security ApiKeyAuth
func AddPermissionToRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = id
	err := models.Find(&role, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	var params AddPermisisonsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}
	for _, permit := range params.Permissions {
		models.Enforcer.AddPermissionForUser(role.AuthKey(), permit)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 角色详情
// @Produce  json
// @Tags 角色权限管理
// @Param id path int true "role id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/roles/{id} [get]
// @Security ApiKeyAuth
func GetRole(c *gin.Context) {
	var role models.Role
	id, _ := strconv.Atoi(c.Param("id"))
	role.ID = id
	err := models.Find(&role, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}
	permissions := role.GetPermissions()
	log.Println("======permission=====", permissions)
	role.Permissions = permissions
	apiHelpers.ResponseSuccess(c, role)
}
