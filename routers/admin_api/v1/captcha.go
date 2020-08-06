package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	config "github.com/xifengzhu/eshop/initializers"
	. "github.com/xifengzhu/eshop/routers/helpers"
)

// @Summary 获取图形验证码
// @Produce  json
// @Tags 后台图形验证码
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/get_captcha [get]
func GetCaptcha(c *gin.Context) {
	id, b64s, err := config.CaptchaGenerate()
	if err != nil {
		ResponseError(c, e.SERVER_ERROR, err.Error())
	}

	data := map[string]interface{}{"img": b64s, "captchaId": id}
	ResponseSuccess(c, data)
}
