package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary 获取微信自定义页面
// @Produce  json
// @Tags 微信页面
// @Param key query string true "page key"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/wxapp_page [get]
func GetWxappPage(c *gin.Context) {
	var wxappPage models.WxappPage
	parmMap := map[string]interface{}{"key": c.Query("key")}
	err := models.Where(Query{Conditions: parmMap}).Find(&wxappPage).Error

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	apiHelpers.ResponseSuccess(c, wxappPage)
}
