package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type PropertyValueParams struct {
	Value          string `json:"value" binding:"required"`
	PropertyNameID int    `json:"property_name_id" binding:"required"`
}

// @Summary 添加规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param params body PropertyValueParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values [post]
// @Security ApiKeyAuth
func AddPropertyValue(c *gin.Context) {
	var err error
	var pvp PropertyValueParams
	if err = c.ShouldBindJSON(&pvp); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var pvalue models.PropertyValue
	copier.Copy(&pvalue, &pvp)

	err = models.SaveResource(&pvalue)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, pvalue)
}

// @Summary 删除规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param id path int true "property_value id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_values/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyValue(c *gin.Context) {
	var pvalue models.PropertyValue
	id, _ := strconv.Atoi(c.Param("id"))
	pvalue.ID = id

	err := models.DestroyResource(pvalue, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}
