package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary 获取分类列表
// @Produce  json
// @Tags 分类
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	models.AllResource(&categories, Query{Preloads: []string{"Children"}})
	response := apiHelpers.Collection{List: categories}
	apiHelpers.ResponseSuccess(c, response)
}
