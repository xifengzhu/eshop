package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 添加主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param params body params.ProductGroupParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/product_groups [post]
// @Security ApiKeyAuth
func AddProductGroup(c *gin.Context) {
	var err error
	var productGroupParams ProductGroupParams
	if err = ValidateParams(c, &productGroupParams, "json"); err != nil {
		return
	}

	var productGroup ProductGroup
	copier.Copy(&productGroup, &productGroupParams)

	err = Save(&productGroup)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, productGroup)
}

// @Summary 删除主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/product_groups/{id} [delete]
// @Security ApiKeyAuth
func DeleteProductGroup(c *gin.Context) {
	var productGroup ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id

	err := Destroy(&productGroup)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 主题产品详情
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/product_groups/{id} [get]
// @Security ApiKeyAuth
func GetProductGroup(c *gin.Context) {
	var productGroup ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id
	err := Find(&productGroup, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	productGroup.Products = productGroup.GetProducts()

	ResponseSuccess(c, productGroup)
}

// @Summary 主题产品列表
// @Produce  json
// @Tags 后台主题产品管理
// @Param params query params.QueryProductGroupParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/product_groups [get]
// @Security ApiKeyAuth
func GetProductGroups(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model ProductGroup
	result := &[]ProductGroup{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 更新主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Param params body params.ProductGroupParams true "productGroup params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/product_groups/{id} [put]
// @Security ApiKeyAuth
func UpdateProductGroup(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
		return
	}
	var err error
	var productGroupParams ProductGroupParams
	if err = ValidateParams(c, &productGroupParams, "json"); err != nil {
		return
	}

	var productGroup ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id
	err = Find(&productGroup, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&productGroup, &productGroupParams)

	err = Save(&productGroup)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, productGroup)
}
