package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"log"
	"strconv"
)

type PropertyParams struct {
	ID             int                   `json:"id,omitempty"`
	Name           string                `json:"name" validate:"required"`
	PropertyValues []PropertyValueParams `json:"property_values"  validate:"dive"`
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
	if err = apiHelpers.ValidateParams(c, &pnp, "json"); err != nil {
		return
	}

	var propertyName models.PropertyName
	copier.Copy(&propertyName, &pnp)

	log.Println("=====add propertyName params===", propertyName)

	err = models.Create(&propertyName)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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

	err := models.DestroyWithCallbacks(propertyName, Query{Callbacks: []func(){propertyName.RemovePropertyValues}})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
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

	err := models.Find(&pn, Query{Preloads: []string{"PropertyValues"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
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

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.PropertyName
	result := &[]models.PropertyName{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q"), Preloads: []string{"PropertyValues"}}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)

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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var propertyParams PropertyParams
	if err = apiHelpers.ValidateParams(c, &propertyParams, "json"); err != nil {
		return
	}

	var propertyName models.PropertyName

	id, _ := strconv.Atoi(c.Param("id"))
	propertyName.ID = id

	err = models.Find(&propertyName, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&propertyName, &propertyParams)
	// reset propertyName goodses
	propertyName.PropertyValues = nil
	// recover the propertyName id
	propertyName.ID = id
	copier.Copy(&propertyName.PropertyValues, &propertyParams.PropertyValues)

	err = propertyName.NestUpdate()
	models.Find(&propertyName, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, propertyName)
}
