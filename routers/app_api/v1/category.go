package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type CategoryProductQueryParams struct {
	utils.Pagination
	CategoryID string `json:"cagegory_id"`
}

// @Summary 获取分类列表
// @Produce  json
// @Tags 分类
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	models.All(&categories, Query{Preloads: []string{"Children"}})
	apiHelpers.ResponseSuccess(c, categories)
}

// @Summary 获取分类商品列表
// @Produce  json
// @Tags 分类
// @Param params query CategoryProductQueryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/categories/{id}/products [get]
func GetCategoryProducts(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var category models.Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id

	models.Find(&category, Query{})

	products := category.GetCategoryProducts(pagination)

	result := transferProductToEntity(products)
	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}
