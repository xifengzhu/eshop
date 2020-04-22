package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

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
