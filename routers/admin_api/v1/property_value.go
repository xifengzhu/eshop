package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 添加规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param params body params.PropertyValueParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_values [post]
// @Security ApiKeyAuth
func AddPropertyValue(c *gin.Context) {
	var err error
	var pvp PropertyValueParams
	if err = ValidateParams(c, &pvp, "json"); err != nil {
		return
	}

	var pvalue models.PropertyValue
	copier.Copy(&pvalue, &pvp)

	err = models.Save(&pvalue)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, pvalue)
}

// @Summary 删除规格值
// @Produce  json
// @Tags 后台规格值管理
// @Param id path int true "property_value id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_values/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyValue(c *gin.Context) {
	var pvalue models.PropertyValue
	id, _ := strconv.Atoi(c.Param("id"))
	pvalue.ID = id

	err := models.Destroy(pvalue)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}
