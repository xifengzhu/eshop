package v1

import (
	"github.com/gin-gonic/gin"
)

type WxappPageParams struct {
}

// @Summary 添加WX页面
// @Produce  json
// @Tags 后台WX页面管理
// @Param params body WxappPageParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages [post]
// @Security ApiKeyAuth
func AddWxappPage(c *gin.Context) {

}

// @Summary 删除WX页面
// @Produce  json
// @Tags 后台WX页面管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [delete]
// @Security ApiKeyAuth
func DeleteWxappPage(c *gin.Context) {

}

// @Summary WX页面详情
// @Produce  json
// @Tags 后台WX页面管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [get]
// @Security ApiKeyAuth
func GetWxappPage(c *gin.Context) {

}

// @Summary 更新WX页面
// @Produce  json
// @Tags 后台WX页面管理
// @Param id path int true "id"
// @Param params body WxappPageParams true "web_page params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [put]
// @Security ApiKeyAuth
func UpdateWxappPage(c *gin.Context) {

}
