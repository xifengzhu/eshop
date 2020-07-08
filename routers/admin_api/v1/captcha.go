package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	config "github.com/xifengzhu/eshop/initializers"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary 获取图形验证码
// @Produce  json
// @Tags 后台图形验证码
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/get_captcha [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := config.CaptchaGenerate()
	if err != nil {
		apiHelpers.ResponseError(c, e.SERVER_ERROR, err.Error())
	}

	data := map[string]interface{}{"img": b64s, "captchaId": id}
	apiHelpers.ResponseSuccess(c, data)
}
