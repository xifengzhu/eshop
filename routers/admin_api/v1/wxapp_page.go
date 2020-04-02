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

type WxAppPageParams struct {
	Name     string `json:"name" binding:"required"`
	PageType int    `json:"page_type" binding:"required"`
	PageData string `json:"page_data" binding:"required"`
}

// @Summary 添加自定义页面
// @Produce  json
// @Tags 后台自定义页面管理
// @Param params body WxAppPageParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages [post]
// @Security ApiKeyAuth
func AddWxAppPage(c *gin.Context) {
	var err error
	var wpParams WxAppPageParams
	if err = c.ShouldBind(&wpParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var wxAppPage models.WxappPage
	copier.Copy(&wxAppPage, &wpParams)

	err = models.SaveResource(&wxAppPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, wxAppPage)
}

// @Summary 删除自定义页面
// @Produce  json
// @Tags 后台自定义页面管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [delete]
// @Security ApiKeyAuth
func DeleteWxAppPage(c *gin.Context) {
	var wxAppPage models.WxappPage
	id, _ := strconv.Atoi(c.Param("id"))
	wxAppPage.ID = id

	err := models.DestroyResource(&wxAppPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 自定义页面详情
// @Produce  json
// @Tags 后台自定义页面管理
// @Param id path int true "web_page id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [get]
// @Security ApiKeyAuth
func GetWxAppPage(c *gin.Context) {
	var wxAppPage models.WxappPage
	id, _ := strconv.Atoi(c.Param("id"))
	wxAppPage.ID = id
	err := models.FindResource(&wxAppPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, wxAppPage)
}

// @Summary 自定义页面列表
// @Produce  json
// @Tags 后台自定义页面管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [get]
// @Security ApiKeyAuth
func GetWxAppPages(c *gin.Context) {
	var wxAppPages []models.WxappPage
	models.AllResource(&wxAppPages, Query{})
	response := apiHelpers.Collection{List: wxAppPages}
	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 更新自定义页面
// @Produce  json
// @Tags 后台自定义页面管理
// @Param id path int true "id"
// @Param params body WxAppPageParams true "web_page params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages/{id} [put]
// @Security ApiKeyAuth
func UpdateWxAppPage(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
		return
	}
	var err error
	var wxAppPageParams WxAppPageParams
	if err = c.ShouldBindJSON(&wxAppPageParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var wxAppPage models.WxappPage
	id, _ := strconv.Atoi(c.Param("id"))
	wxAppPage.ID = id
	err = models.FindResource(&wxAppPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&wxAppPage, &wxAppPageParams)

	err = models.SaveResource(&wxAppPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, wxAppPage)
}
