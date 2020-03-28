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
	title := c.Query("title")
	var webPage models.WebPage
	err := webPage.FindByTitle(title)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, webPage)
}
