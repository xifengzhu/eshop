package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary 获取webview页面
// @Produce  json
// @Tags 微信页面
// @Param title query string true "page title"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/web_page [get]
func GetWebPage(c *gin.Context) {
	var webPage models.WebPage
	parmMap := map[string]interface{}{"title": c.Query("title")}
	err := models.FirstResource(&webPage, Query{Conditions: parmMap})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, webPage)
}
