package v1

import (
	"errors"
	// "fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"

	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

type AdminUserParams struct {
	Email           string `json:"email"  binding:"required,email"`
	Password        string `json:"password"  binding:"required,eqfield=ConfirmPassword,gte=6,lt=12"`
	ConfirmPassword string `json:"confirm_password"  binding:"required"`
}

type LoginParams struct {
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type QueryAdminUserParams struct {
	utils.Pagination
	Email_cont      string    `json:"q[email_cont]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Status          []int     `json:"q[status_in]"`
}

// @Summary 管理员登录
// @Produce  json
// @Tags 后台管理员
// @Param params body LoginParams true "邮箱密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/login [post]
func Login(c *gin.Context) {
	var login LoginParams
	if err := c.ShouldBindJSON(&login); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var admin models.AdminUser
	err := admin.GetAdminUserByEmail(login.Email)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("账号不存在"))
		return
	}

	if admin.Authenticate(login.Password) {
		var params = map[string]interface{}{"id": admin.ID, "resource": "admin"}
		token := utils.Encode(params)
		apiHelpers.ResponseSuccess(c, token)
		return
	}
	apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("密码错误"))
}

// @Summary 添加管理员
// @Produce  json
// @Tags 后台管理员
// @Param params body AdminUserParams true "邮箱密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/admin_users [post]
// @Security ApiKeyAuth
func AddAdminUser(c *gin.Context) {
	var err error
	var adminUserParams AdminUserParams
	if err = c.ShouldBindJSON(&adminUserParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var admin models.AdminUser
	err = admin.GetAdminUserByEmail(adminUserParams.Email)
	if err == nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("账号已经存在"))
		return
	}

	var adminUser = models.AdminUser{Email: adminUserParams.Email, Password: adminUserParams.Password}
	err = adminUser.Create()
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("添加账号失败"))
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

	models.SearchResourceQuery(&model, result, &pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}
