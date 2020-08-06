package v1

import (
	"github.com/gin-gonic/gin"
	. "github.com/xifengzhu/eshop/models"
	_ "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 获取分类列表
// @Produce  json
// @Tags 分类
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/categories [get]
func GetCategories(c *gin.Context) {
	var categories []Category
	All(&categories, Options{Preloads: []string{"Children"}})
	ResponseSuccess(c, categories)
}

// @Summary 获取分类商品列表
// @Produce  json
// @Tags 分类
// @Param params query params.CategoryProductQueryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/categories/{id}/products [get]
func GetCategoryProducts(c *gin.Context) {
	pagination := SetDefaultPagination(c)

	var category Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id

	Find(&category, Options{})

	products := category.GetCategoryProducts(pagination)

	result := transferProductToEntity(products)
	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}
