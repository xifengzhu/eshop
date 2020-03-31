package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type PropertyParams struct {
	Name           string   `json:"name" binding:"required"`
	PropertyValues []string `json:"property_values"`
}

type QueryPropertyParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}

// @Summary 添加规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param params body PropertyParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names [post]
// @Security ApiKeyAuth
func AddPropertyName(c *gin.Context) {
	var err error
	var pnp PropertyParams
	if err = c.ShouldBindJSON(&pnp); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var propertyName models.PropertyName
	copier.Copy(&propertyName, &pnp)

	err = models.SaveResource(&propertyName)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, propertyName)
}

// @Summary 删除规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyName(c *gin.Context) {
	var propertyName models.PropertyName
	id, _ := strconv.Atoi(c.Param("id"))
	propertyName.ID = id

	err := models.DestroyResource(propertyName, Query{Callbacks: []func(){propertyName.RemovePropertyValues}})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 规格名详情
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [get]
// @Security ApiKeyAuth
func GetPropertyName(c *gin.Context) {
	var pn models.PropertyName
	id, _ := strconv.Atoi(c.Param("id"))
	pn.ID = id

	err := models.FindResource(&pn, Query{Preloads: []string{"PropertyValues"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, pn)
}

// @Summary 规格列表
// @Produce  json
// @Tags 后台规格名管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names [get]
// @Security ApiKeyAuth
func GetPropertyNames(c *gin.Context) {
	var pnames []models.PropertyName
	models.AllResource(&pnames, Query{Preloads: []string{"PropertyValues"}})
	apiHelpers.ResponseSuccess(c, pnames)
}

// @Summary 更新规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "id"
// @Param params body PropertyParams true "property_name params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/property_names/{id} [put]
// @Security ApiKeyAuth
func UpdatePropertyName(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
	}
	var err error
	var propertyParams PropertyParams
	if err = c.ShouldBindJSON(&propertyParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var propertyName models.PropertyName

	id, _ := strconv.Atoi(c.Param("id"))
	propertyName.ID = id

	err = models.FindResource(&propertyName, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&propertyName, &propertyParams)
	// reset propertyName goodses
	propertyName.PropertyValues = nil
	// recover the propertyName id
	propertyName.ID = id
	copier.Copy(&propertyName.PropertyValues, &propertyParams.PropertyValues)

	err = propertyName.NestUpdate()
	models.FindResource(&propertyName, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, propertyName)
}
