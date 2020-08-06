package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 添加管理员
// @Produce  json
// @Tags 后台管理员
// @Param params body params.AddAdminUserParams true "邮箱密码"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users [post]
// @Security ApiKeyAuth
func AddAdminUser(c *gin.Context) {
	var err error
	var adminUserParams AddAdminUserParams
	if err = c.ShouldBindJSON(&adminUserParams); err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var admin AdminUser
	parmMap := map[string]interface{}{"email": adminUserParams.Email}
	exist := Exist(&admin, Options{Conditions: parmMap})
	if exist {
		ResponseError(c, e.INVALID_PARAMS, "账号已经存在")
		return
	}

	copier.Copy(&admin, &adminUserParams)
	admin.Status = "active"
	err = Create(&admin)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, "添加账号失败")
		return
	}

	ResponseOK(c)
}

// @Summary 更新管理员
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin user id"
// @Param params body params.UpdateAdminUserParams true "邮箱密码"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users/{id} [put]
// @Security ApiKeyAuth
func UpdateAdminUser(c *gin.Context) {
	var err error
	var adminUserParams UpdateAdminUserParams
	if err = c.ShouldBindJSON(&adminUserParams); err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var adminUser AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err = Find(&adminUser, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&adminUser, &adminUserParams)
	err = Save(&adminUser)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, "更新账号失败")
		return
	}

	ResponseOK(c)
}

// @Summary 管理员列表
// @Produce  json
// @Tags 后台管理员
// @Param params query params.QueryAdminUserParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users [get]
// @Security ApiKeyAuth
func GetAdminUsers(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model AdminUser
	result := &[]AdminUser{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 管理员详情
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin user id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users/{id} [get]
// @Security ApiKeyAuth
func GetAdminUser(c *gin.Context) {
	var adminUser AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err := Find(&adminUser, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	adminUser.Roles = adminUser.GetRoles()
	adminUser.Permissions = adminUser.GetPermissions()

	ResponseSuccess(c, adminUser)
}

// @Summary 删除管理员
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin_user id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users/{id} [delete]
// @Security ApiKeyAuth
func DeleteAdminUser(c *gin.Context) {
	var admin AdminUser
	admin.ID, _ = strconv.Atoi(c.Param("id"))

	err := Destroy(&admin)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseOK(c)
}

// @Summary 分配角色给用户
// @Produce  json
// @Tags 后台管理员
// @Param params body params.AddRolesParams true "角色ID数组"
// @Param id path int true "管理员id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_users/{id}/roles [post]
// @Security ApiKeyAuth
func AddRoleForUser(c *gin.Context) {
	var adminUser AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err := Find(&adminUser, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	var params AddRolesParams
	if err = c.ShouldBindJSON(&params); err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var roles []Role
	Where(Options{Conditions: params.RoleIds}).Find(&roles)

	for _, role := range roles {
		Enforcer.AddRoleForUser(adminUser.AuthKey(), role.Name)
	}
	ResponseOK(c)
}

// @Summary 当前管理员权限
// @Produce  json
// @Tags 后台管理员
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/admin_user/abilities [get]
// @Security ApiKeyAuth
func GetAdminUserPermissions(c *gin.Context) {
	adminStr, _ := c.Get("resource")
	adminUser := adminStr.(AdminUser)
	permissions := adminUser.GetPermissions()
	ResponseSuccess(c, permissions)
}
