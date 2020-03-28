package v1

import (
	"github.com/gin-gonic/gin"
)

type ProductParams struct {
}

// @Summary 添加产品
// @Produce  json
// @Tags 后台产品管理
// @Param params body ProductParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products [post]
// @Security ApiKeyAuth
func AddProduct(c *gin.Context) {

}

// @Summary 删除产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [delete]
// @Security ApiKeyAuth
func DeleteProduct(c *gin.Context) {

}

// @Summary 产品详情
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [get]
// @Security ApiKeyAuth
func GetProduct(c *gin.Context) {

}

// @Summary 更新产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "id"
// @Param params body ProductParams true "product params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [put]
// @Security ApiKeyAuth
func UpdateProduct(c *gin.Context) {

}
