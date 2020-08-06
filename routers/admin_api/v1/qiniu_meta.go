package v1

import (
	"github.com/gin-gonic/gin"
	config "github.com/xifengzhu/eshop/initializers"
	apiHelpers "github.com/xifengzhu/eshop/routers/helpers"
)

// @Summary 七牛上传凭证
// @Produce  json
// @Tags 上传凭证
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/qiniu_meta [get]
// @Security ApiKeyAuth
func GetQiniuMeta(c *gin.Context) {
	meta := config.GetUploadMeta()
	apiHelpers.ResponseSuccess(c, meta)
}
