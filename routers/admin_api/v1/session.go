package v1

import (
	"errors"
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

type CaptchaParam struct {
	CaptchaID    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

type LoginParams struct {
	CaptchaParam
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password"  binding:"required"`
}

type ForgetPasswordParams struct {
	CaptchaParam
	Email string `json:"email"  binding:"required,email"`
}

type ResetPasswordParams struct {
	Password string `json:"password"  binding:"required,gte=6,lt=12"`
}

// @Summary 管理员登录
// @Produce  json
// @Tags 后台管理员
// @Param params body LoginParams true "邮箱密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/sessions/login [post]
func Login(c *gin.Context) {
	var login LoginParams
	if err := c.ShouldBindJSON(&login); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	valid := utils.CaptchaVerify(login.CaptchaID, login.CaptchaValue)
	if !valid {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("验证码错误"))
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

// @Summary 当前管理员
// @Produce  json
// @Tags 后台管理员
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/sessions/mine [get]
// @Security ApiKeyAuth
func GetCurrentAdminUsers(c *gin.Context) {
	admin, _ := c.Get("resource")
	apiHelpers.ResponseSuccess(c, admin.(models.AdminUser))
}

// @Summary 管理员忘记密码
// @Produce  json
// @Tags 后台管理员
// @Param params body ForgetPasswordParams true "重置密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/sessions/forget_password [post]
func ForgetPassword(c *gin.Context) {
	var resetParam ForgetPasswordParams
	if err := c.ShouldBindJSON(&resetParam); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	valid := utils.CaptchaVerify(resetParam.CaptchaID, resetParam.CaptchaValue)
	if !valid {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("验证码错误"))
		return
	}

	var admin models.AdminUser
	err := admin.GetAdminUserByEmail(resetParam.Email)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("账号不存在"))
		return
	}
	// TODO: send email

}

// @Summary 管理员设置密码
// @Produce  json
// @Tags 后台管理员
// @Param reset_password_token path string true "充值密码token"
// @Param params body ResetPasswordParams true "重置密码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/sessions/reset_password [put]
func ResetPassword(c *gin.Context) {

}
