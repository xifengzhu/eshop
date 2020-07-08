package v1

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

type AddAdminUserParams struct {
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password"  validate:"required,gte=6,lt=12"`
}

type UpdateAdminUserParams struct {
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Status   string `json:"status,omitempty" validate:"oneof=active banned"`
}

type AddRolesParams struct {
	RoleIds []int `json:"role_ids"`
}

type QueryAdminUserParams struct {
	utils.Pagination
	Email_cont      string    `json:"q[email_cont]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Status          []int     `json:"q[status_in]"`
}

// @Summary 添加管理员
// @Produce  json
// @Tags 后台管理员
// @Param params body AddAdminUserParams true "邮箱密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users [post]
// @Security ApiKeyAuth
func AddAdminUser(c *gin.Context) {
	var err error
	var adminUserParams AddAdminUserParams
	if err = c.ShouldBindJSON(&adminUserParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var admin models.AdminUser
	parmMap := map[string]interface{}{"email": adminUserParams.Email}
	exist := models.Exist(&admin, Query{Conditions: parmMap})
	if exist {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "账号已经存在")
		return
	}

	copier.Copy(&admin, &adminUserParams)
	admin.Status = "active"
	err = models.Create(&admin)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "添加账号失败")
		return
	}

	apiHelpers.ResponseOK(c)
}

// @Summary 更新管理员
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin user id"
// @Param params body UpdateAdminUserParams true "邮箱密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users/{id} [put]
// @Security ApiKeyAuth
func UpdateAdminUser(c *gin.Context) {
	var err error
	var adminUserParams UpdateAdminUserParams
	if err = c.ShouldBindJSON(&adminUserParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var adminUser models.AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err = models.Find(&adminUser, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&adminUser, &adminUserParams)
	err = models.Save(&adminUser)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "更新账号失败")
		return
	}

	apiHelpers.ResponseOK(c)
}

// @Summary 管理员列表
// @Produce  json
// @Tags 后台管理员
// @Param params query QueryAdminUserParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users [get]
// @Security ApiKeyAuth
func GetAdminUsers(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.AdminUser
	result := &[]models.AdminUser{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 管理员详情
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin user id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users/{id} [get]
// @Security ApiKeyAuth
func GetAdminUser(c *gin.Context) {
	var adminUser models.AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err := models.Find(&adminUser, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	adminUser.Roles = adminUser.GetRoles()
	adminUser.Permissions = adminUser.GetPermissions()

	apiHelpers.ResponseSuccess(c, adminUser)
}

// @Summary 删除管理员
// @Produce  json
// @Tags 后台管理员
// @Param id path int true "admin_user id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users/{id} [delete]
// @Security ApiKeyAuth
func DeleteAdminUser(c *gin.Context) {
	var admin models.AdminUser
	admin.ID, _ = strconv.Atoi(c.Param("id"))

	err := models.Destroy(&admin)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 分配角色给用户
// @Produce  json
// @Tags 后台管理员
// @Param params body AddRolesParams true "角色ID数组"
// @Param id path int true "管理员id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users/{id}/roles [post]
// @Security ApiKeyAuth
func AddRoleForUser(c *gin.Context) {
	var adminUser models.AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err := models.Find(&adminUser, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	var params AddRolesParams
	if err = c.ShouldBindJSON(&params); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	var roles []models.Role
	models.Where(Query{Conditions: params.RoleIds}).Find(&roles)

	for _, role := range roles {
		models.Enforcer.AddRoleForUser(adminUser.AuthKey(), role.Name)
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 当前管理员权限
// @Produce  json
// @Tags 后台管理员
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_user/abilities [get]
// @Security ApiKeyAuth
func GetAdminUserPermissions(c *gin.Context) {
	adminStr, _ := c.Get("resource")
	adminUser := adminStr.(models.AdminUser)
	permissions := adminUser.GetPermissions()
	apiHelpers.ResponseSuccess(c, permissions)
}
