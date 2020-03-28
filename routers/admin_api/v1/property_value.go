package v1

import (
	"github.com/gin-gonic/gin"
)

type PropertyValueParams struct {
}

// @Summary 添加规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param params body PropertyValueParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values [post]
// @Security ApiKeyAuth
func AddPropertyValue(c *gin.Context) {

}

// @Summary 删除规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param id path int true "property_value id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyValue(c *gin.Context) {

}

// @Summary 规格值详情
// @Produce  json
// @Tags 后台规格值管理
// @Param id path int true "property_value id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values/{id} [get]
// @Security ApiKeyAuth
func GetPropertyValue(c *gin.Context) {

}

// @Summary 更新规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param id path int true "id"
// @Param params body PropertyValueParams true "property_value params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values/{id} [put]
// @Security ApiKeyAuth
func UpdatePropertyValue(c *gin.Context) {

}
