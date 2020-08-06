package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	config "github.com/xifengzhu/eshop/initializers"
	"github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
)

// @Summary 管理员登录
// @Produce  json
// @Tags 后台管理员
// @Param params body params.LoginParams true "邮箱密码"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/sessions/login [post]
func Login(c *gin.Context) {
	var login LoginParams

	if err := ValidateParams(c, &login, "json"); err != nil {
		return
	}

	valid := config.CaptchaVerify(login.CaptchaID, login.CaptchaValue)
	if !valid {
		ResponseError(c, e.INVALID_PARAMS, "验证码错误")
		return
	}

	var admin models.AdminUser
	err := admin.GetAdminUserByEmail(login.Email)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, "账号不存在")
		return
	}

	if admin.Authenticate(login.Password) {
		var params = map[string]interface{}{"id": admin.ID, "resource": "admin"}
		token := utils.Encode(params)
		ResponseSuccess(c, token)
		return
	}
	ResponseError(c, e.INVALID_PARAMS, "密码错误")
}

// @Summary 当前管理员
// @Produce  json
// @Tags 后台管理员
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/sessions/mine [get]
// @Security ApiKeyAuth
func GetCurrentAdminUser(c *gin.Context) {
	admin, _ := c.Get("resource")
	ResponseSuccess(c, admin.(models.AdminUser))
}

// @Summary 管理员忘记密码
// @Produce  json
// @Tags 后台管理员
// @Param params body params.ForgetPasswordParams true "重置密码"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/sessions/forget_password [post]
func ForgetPassword(c *gin.Context) {
	var resetParam ForgetPasswordParams
	if err := ValidateParams(c, &resetParam, "json"); err != nil {
		return
	}

	valid := config.CaptchaVerify(resetParam.CaptchaID, resetParam.CaptchaValue)
	if !valid {
		ResponseError(c, e.INVALID_PARAMS, "验证码错误")
		return
	}

	var admin models.AdminUser
	err := admin.GetAdminUserByEmail(resetParam.Email)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, "账号不存在")
		return
	}
	// TODO: send email

}

// @Summary 管理员设置密码
// @Produce  json
// @Tags 后台管理员
// @Param reset_password_token path string true "充值密码token"
// @Param params body params.ResetPasswordParams true "重置密码"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/sessions/reset_password [put]
func ResetPassword(c *gin.Context) {

}
