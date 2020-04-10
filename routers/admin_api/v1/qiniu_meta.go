package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/utils"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary 七牛上传凭证
// @Produce  json
// @Tags 上传凭证
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/qiniu_meta [get]
// @Security ApiKeyAuth
func GetQiniuMeta(c *gin.Context) {
	meta := utils.GetUploadMeta()
	apiHelpers.ResponseSuccess(c, meta)
}
