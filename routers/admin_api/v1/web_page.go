package v1

import (
	"github.com/gin-gonic/gin"
)

type WebPageParams struct {
}

// @Summary 添加webview
// @Produce  json
// @Tags 后台webview管理
// @Param params body WebPageParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages [post]
// @Security ApiKeyAuth
func AddWebPage(c *gin.Context) {

}

// @Summary 删除webview
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages/{id} [delete]
// @Security ApiKeyAuth
func DeleteWebPage(c *gin.Context) {

}

// @Summary webview详情
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages/{id} [get]
// @Security ApiKeyAuth
func GetWebPage(c *gin.Context) {

}

// @Summary 更新webview
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "id"
// @Param params body WebPageParams true "web_page params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages/{id} [put]
// @Security ApiKeyAuth
func UpdateWebPage(c *gin.Context) {

}
