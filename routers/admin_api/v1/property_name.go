package v1

import (
	"github.com/gin-gonic/gin"
)

type PropertyNameParams struct {
}

// @Summary 添加规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param params body PropertyNameParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names [post]
// @Security ApiKeyAuth
func AddPropertyName(c *gin.Context) {

}

// @Summary 删除规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyName(c *gin.Context) {

}

// @Summary 规格名详情
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [get]
// @Security ApiKeyAuth
func GetPropertyName(c *gin.Context) {

}

// @Summary 更新规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "id"
// @Param params body PropertyNameParams true "property_name params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [put]
// @Security ApiKeyAuth
func UpdatePropertyName(c *gin.Context) {

}
