package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

type CaptchaParam struct {
	ID    string `json:"id" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// @Summary 获取图形验证码
// @Produce  json
// @Tags 后台图形验证码
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/get_captcha [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := utils.CaptchaGenerate()
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR, err)
	}

	data := map[string]interface{}{"img": b64s, "captchaId": id}
	apiHelpers.ResponseSuccess(c, data)
}

// @Summary 校验图形验证码
// @Produce  json
// @Tags 后台图形验证码
// @Param params body CaptchaParam true "验证码"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/verify_captcha [post]
func VerifyCaptcha(c *gin.Context) {
	var err error
	var captchaParam CaptchaParam
	if err = c.ShouldBindJSON(&captchaParam); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}

	valid := utils.CaptchaVerify(captchaParam.ID, captchaParam.Value)
	apiHelpers.ResponseSuccess(c, valid)
}
