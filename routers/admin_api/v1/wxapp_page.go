package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type WxAppPageParams struct {
	Name                 string      `json:"name,omitempty"`
	Key                  string      `json:"key,omitempty"`
	PageData             models.JSON `json:"page_data,omitempty"`
	ShareSentence        string      `json:"share_sentence,omitempty"`
	ShareCover           string      `json:"share_cover,omitempty"`
	ShareBackgroundCover string      `json:"share_background_cover,omitempty"`
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
	if err = apiHelpers.ValidateParams(c, &wpParams, "json"); err != nil {
		return
	}

	var wxAppPage models.WxappPage
	copier.Copy(&wxAppPage, &wpParams)
	wxAppPage.PageType = "2"

	err = models.Save(&wxAppPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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

	err := models.Destroy(&wxAppPage)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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
	err := models.Find(&wxAppPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, wxAppPage)
}

// @Summary 自定义页面列表
// @Produce  json
// @Tags 后台自定义页面管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxapp_pages [get]
// @Security ApiKeyAuth
func GetWxAppPages(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.WxappPage
	result := &[]models.WxappPage{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
		return
	}
	var err error
	var wxAppPageParams WxAppPageParams
	if err = apiHelpers.ValidateParams(c, &wxAppPageParams, "json"); err != nil {
		return
	}

	var wxAppPage models.WxappPage
	id, _ := strconv.Atoi(c.Param("id"))
	wxAppPage.ID = id
	err = models.Find(&wxAppPage, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	changedAttrs := models.WxappPage{}
	copier.Copy(&changedAttrs, &wxAppPageParams)

	err = models.Update(&wxAppPage, changedAttrs)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, wxAppPage)
}

func GetPageGroupLinks(c *gin.Context) {
	pageLinks := map[string]interface{}{
		"列表": map[string]string{
			"产品列表": "/extra/products/pages/list",
		},
		"单个产品": map[string]string{
			"单骰重摇印花T恤(白色)": "/extra/products/pages/detail?id=88",
			"单骰重摇印花T恤(黑色)": "/extra/products/pages/detail?id=89",
		},
		"新品/折扣/推荐": map[string]string{
			"周一新品":  "/extra/shares/pages/custom-products-list?tab_module_id=1&title=周一新品",
			"折扣商品":  "/extra/shares/pages/custom-products-list?tab_module_id=2&title=折扣商品",
			"推荐商品":  "/extra/shares/pages/custom-products-list?tab_module_id=3&title=推荐商品",
			"纯色基础款": "/extra/shares/pages/custom-products-list?tab_module_id=4&title=纯色基础款",
		},
		"普通分类": map[string]string{
			"T恤": "/extra/products/pages/list?category_id=6",
			"上装": "/extra/products/pages/list?category_id=1",
		},
		"独立页面": map[string]string{
			"关于我们": "/extra/shares/pages/webview?key=about_us",
			"尺码助手": "/extra/shares/pages/webview?key=size_description",
			"常见问题": "/extra/shares/pages/webview?key=qa",
			"退货政策": "/extra/shares/pages/webview?key=return_policy",
		},
		"自定义页面": map[string]string{
			"品牌故事":  "/extra/shares/pages/custom-page?enName=Brandstory&title=品牌故事",
			"热门分类":  "/extra/shares/pages/custom-page?enName=hot_category&title=热门分类",
			"纯色基础款": "/extra/shares/pages/custom-page?enName=basistee&title=纯色基础款",
			"首页区块":  "/extra/shares/pages/custom-page?enName=home&title=首页区块",
		},
	}
	apiHelpers.ResponseSuccess(c, pageLinks)
}
