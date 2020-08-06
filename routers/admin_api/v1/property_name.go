package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"log"
	"strconv"
)

// @Summary 添加规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param params body params.PropertyParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_names [post]
// @Security ApiKeyAuth
func AddPropertyName(c *gin.Context) {
	var err error
	var pnp PropertyParams
	if err = ValidateParams(c, &pnp, "json"); err != nil {
		return
	}

	var propertyName PropertyName
	copier.Copy(&propertyName, &pnp)

	log.Println("=====add propertyName params===", propertyName)

	err = Create(&propertyName)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, propertyName)
}

// @Summary 删除规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_names/{id} [delete]
// @Security ApiKeyAuth
func DeletePropertyName(c *gin.Context) {
	var propertyName PropertyName
	id, _ := strconv.Atoi(c.Param("id"))
	propertyName.ID = id

	err := DestroyWithCallbacks(propertyName, Options{Callbacks: []func(){propertyName.RemovePropertyValues}})
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 规格名详情
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "property_name id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_names/{id} [get]
// @Security ApiKeyAuth
func GetPropertyName(c *gin.Context) {
	var pn PropertyName
	id, _ := strconv.Atoi(c.Param("id"))
	pn.ID = id

	err := Find(&pn, Options{Preloads: []string{"PropertyValues"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, pn)
}

// @Summary 规格列表
// @Produce  json
// @Tags 后台规格名管理
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_names [get]
// @Security ApiKeyAuth
func GetPropertyNames(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model PropertyName
	result := &[]PropertyName{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q"), Preloads: []string{"PropertyValues"}}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)

}

// @Summary 更新规格名
// @Produce  json
// @Tags 后台规格名管理
// @Param id path int true "id"
// @Param params body params.PropertyParams true "property_name params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/property_names/{id} [put]
// @Security ApiKeyAuth
func UpdatePropertyName(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var propertyParams PropertyParams
	if err = ValidateParams(c, &propertyParams, "json"); err != nil {
		return
	}

	var propertyName PropertyName

	id, _ := strconv.Atoi(c.Param("id"))
	propertyName.ID = id

	err = Find(&propertyName, Options{})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&propertyName, &propertyParams)
	// reset propertyName goodses
	propertyName.PropertyValues = nil
	// recover the propertyName id
	propertyName.ID = id
	copier.Copy(&propertyName.PropertyValues, &propertyParams.PropertyValues)

	err = propertyName.NestUpdate()
	Find(&propertyName, Options{})
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, propertyName)
}
