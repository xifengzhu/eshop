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
// @Param title query string true "page title"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/wxapp_page [get]
func GetWxappPage(c *gin.Context) {
	var wxappPage models.WxappPage
	parmMap := map[string]interface{}{"name": c.Query("name")}
	err := models.FirstResource(&wxappPage, Query{Conditions: parmMap})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, wxappPage)
}
