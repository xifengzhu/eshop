package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
	// "time"
)

// @Summary 添加产品
// @Produce  json
// @Tags 后台产品管理
// @Param params body params.ProductParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products [post]
// @Security ApiKeyAuth
func AddProduct(c *gin.Context) {
	var err error
	var productParams ProductParams
	if err = ValidateParams(c, &productParams, "json"); err != nil {
		return
	}

	var product Product
	copier.Copy(&product, &productParams)

	err = Save(&product)

	// update categories
	var categories []Category
	Where(Options{Conditions: productParams.CategoryIDs}).Find(&categories)
	fmt.Println("========categories=======", categories)

	product.UpdateCatgories(categories)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, product)
}

// @Summary 删除产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products/{id} [delete]
// @Security ApiKeyAuth
func DeleteProduct(c *gin.Context) {
	var product Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := DestroyWithCallbacks(product, Options{Callbacks: []func(){product.RemoveGoodses}})
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 产品列表
// @Produce  json
// @Tags 后台产品管理
// @Param params query params.QueryProductParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products [get]
// @Security ApiKeyAuth
func GetProductes(c *gin.Context) {
	pagination := SetDefaultPagination(c)

	var model Product
	result := &[]Product{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 产品详情
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products/{id} [get]
// @Security ApiKeyAuth
func GetProduct(c *gin.Context) {
	var product Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := Find(&product, Options{Preloads: []string{"Goodses", "Delivery", "Categories"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, product)
}

// @Summary SKU列表
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product_id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products/{id}/goodses [get]
// @Security ApiKeyAuth
func GetGoodses(c *gin.Context) {
	var goodses []Goods
	fmt.Println("=======id======", c.Param("id"))
	parmMap := map[string]interface{}{"product_id": c.Param("id")}
	Where(Options{Conditions: parmMap}).Find(&goodses)
	ResponseSuccess(c, goodses)
}

// @Summary 更新产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "id"
// @Param params body params.ProductParams true "product params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/products/{id} [put]
// @Security ApiKeyAuth
func UpdateProduct(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var productParams ProductParams
	if err = ValidateParams(c, &productParams, "json"); err != nil {
		return
	}

	fmt.Println("======productParams======", productParams)

	var product Product

	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err = Find(&product, Options{})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&product, &productParams)
	// reset product goodses
	product.Goodses = nil
	// recover the product id
	product.ID = id
	copier.Copy(&product.Goodses, &productParams.Goodses)

	err = product.NestUpdate()

	var categories []Category
	Where(Options{Conditions: productParams.CategoryIDs}).Find(&categories)
	fmt.Println("========categories=======", categories)
	product.UpdateCatgories(categories)

	Find(&product, Options{})
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, product)
}
