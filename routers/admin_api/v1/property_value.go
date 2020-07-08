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
	ID             int    `json:"id,omitempty"`
	Value          string `json:"value" validate:"required"`
	PropertyNameID int    `json:"property_name_id" validate:"required"`
	Destroy        bool   `json:"_destroy,omitempty"`
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
	if err = apiHelpers.ValidateParams(c, &pvp); err != nil {
		return
	}

	var pvalue models.PropertyValue
	copier.Copy(&pvalue, &pvp)

	err = models.Save(&pvalue)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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

	err := models.Destroy(pvalue)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}
