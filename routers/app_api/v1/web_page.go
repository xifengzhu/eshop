package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/helpers"

	"log"
)

// @Summary 获取webview页面
// @Produce  json
// @Tags 微信页面
// @Param title query string true "page title"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/web_page [get]
func GetWebPage(c *gin.Context) {
	var webPage WebPage
	log.Println("===query===", c.Query("title"))
	parmMap := map[string]interface{}{"title": c.Query("title")}
	err := Find(&webPage, Options{Conditions: parmMap})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, webPage)
}
