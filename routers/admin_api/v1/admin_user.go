package v1

import (
	"errors"
	// "fmt"
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
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password"  binding:"required,gte=6,lt=12"`
}

type UpdateAdminUserParams struct {
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Status   string `json:"status,omitempty" binding:"oneof=active banned"`
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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var admin models.AdminUser
	parmMap := map[string]interface{}{"email": adminUserParams.Email}
	exist := models.ExistResource(&admin, Query{Conditions: parmMap})
	if exist {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("账号已经存在"))
		return
	}

	copier.Copy(&admin, &adminUserParams)
	admin.Status = "active"
	err = models.CreateResource(&admin)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("添加账号失败"))
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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var adminUser models.AdminUser
	id, _ := strconv.Atoi(c.Param("id"))
	adminUser.ID = id
	err = models.FindResource(&adminUser, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&adminUser, &adminUserParams)
	err = models.SaveResource(&adminUser)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("更新账号失败"))
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

	models.SearchResourceQuery(&model, result, pagination, c.QueryMap("q"))

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
	err := models.FindResource(&adminUser, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

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

	err := models.DestroyResource(&admin, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseOK(c)
}
