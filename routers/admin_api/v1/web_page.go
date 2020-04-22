package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type WebPageParams struct {
	Title   string `json:"title" binding:"required"`
	Cover   string `json:"cover"`
	Content string `json:"content" binding:"required"`
}

// @Summary 添加webview
// @Produce  json
// @Tags 后台webview管理
// @Param params body WebPageParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages [post]
// @Security ApiKeyAuth
func AddWebPage(c *gin.Context) {
	var err error
	var wpParams WebPageParams
	if err = c.ShouldBind(&wpParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var webPage models.WebPage
	copier.Copy(&webPage, &wpParams)

	err = models.SaveResource(&webPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, webPage)
}

// @Summary 删除webview
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages/{id} [delete]
// @Security ApiKeyAuth
func DeleteWebPage(c *gin.Context) {
	var webPage models.WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id

	err := models.DestroyResource(&webPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary webview详情
// @Produce  json
// @Tags 后台webview管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages/{id} [get]
// @Security ApiKeyAuth
func GetWebPage(c *gin.Context) {
	var webPage models.WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id
	err := models.FindResource(&webPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, webPage)
}

// @Summary webview列表
// @Produce  json
// @Tags 后台webview管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/web_pages [get]
// @Security ApiKeyAuth
func GetWebPages(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.WebPage
	result := &[]models.WebPage{}

	models.SearchResourceQuery(&model, result, pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
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
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
		return
	}
	var err error
	var webPageParams WebPageParams
	if err = c.ShouldBindJSON(&webPageParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var webPage models.WebPage
	id, _ := strconv.Atoi(c.Param("id"))
	webPage.ID = id
	err = models.FindResource(&webPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&webPage, &webPageParams)

	err = models.SaveResource(&webPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, webPage)
}
