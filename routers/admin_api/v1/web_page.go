package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 添加webview
// @Produce  json
// @Tags 后台webview管理
// @Param params body params.WebPageParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/web_pages [post]
// @Security ApiKeyAuth
func AddWebPage(c *gin.Context) {
	var err error
	var wpParams WebPageParams
	if err = ValidateParams(c, &wpParams, "json"); err != nil {
		return
	}

	var webPage WebPage
	copier.Copy(&webPage, &wpParams)

	err = Save(&webPage)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, webPage)
}

// @Summary 删除webview
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/web_pages/{id} [delete]
// @Security ApiKeyAuth
func DeleteWebPage(c *gin.Context) {
	var webPage WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id

	err := Destroy(&webPage)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary webview详情
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/web_pages/{id} [get]
// @Security ApiKeyAuth
func GetWebPage(c *gin.Context) {
	var webPage WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id
	err := Find(&webPage, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, webPage)
}

// @Summary webview列表
// @Produce  json
// @Tags 后台webview管理
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/web_pages [get]
// @Security ApiKeyAuth
func GetWebPages(c *gin.Context) {
	pagination := SetDefaultPagination(c)

	var model WebPage
	result := &[]WebPage{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 更新webview
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "id"
// @Param params body params.WebPageParams true "web_page params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/web_pages/{id} [put]
// @Security ApiKeyAuth
func UpdateWebPage(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
		return
	}
	var err error
	var webPageParams WebPageParams
	if err := ValidateParams(c, &webPageParams, "json"); err != nil {
		return
	}

	var webPage WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id
	err = Find(&webPage, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&webPage, &webPageParams)

	err = Save(&webPage)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, webPage)
}
