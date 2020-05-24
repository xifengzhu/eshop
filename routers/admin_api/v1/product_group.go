package v1

import (
	"errors"
	// "fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
	"time"
)

type ProductGroupParams struct {
	Name       string      `json:"name"`
	Remark     string      `json:"remark"`
	ProductIDs models.JSON `json:"product_ids"`
	Key        string      `json:"key"`
}

type QueryProductGroupParams struct {
	utils.Pagination
	Name            string    `json:"q[name]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
}

// @Summary 添加主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param params body ProductGroupParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/product_groups [post]
// @Security ApiKeyAuth
func AddProductGroup(c *gin.Context) {
	var err error
	var productGroupParams ProductGroupParams
	if err = c.ShouldBind(&productGroupParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var productGroup models.ProductGroup
	copier.Copy(&productGroup, &productGroupParams)

	err = models.Save(&productGroup)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, productGroup)
}

// @Summary 删除主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/product_groups/{id} [delete]
// @Security ApiKeyAuth
func DeleteProductGroup(c *gin.Context) {
	var productGroup models.ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id

	err := models.Destroy(&productGroup)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 主题产品详情
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/product_groups/{id} [get]
// @Security ApiKeyAuth
func GetProductGroup(c *gin.Context) {
	var productGroup models.ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id
	err := models.Find(&productGroup, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	productGroup.Products = productGroup.GetProducts()

	apiHelpers.ResponseSuccess(c, productGroup)
}

// @Summary 主题产品列表
// @Produce  json
// @Tags 后台主题产品管理
// @Param params query QueryProductGroupParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/product_groups [get]
// @Security ApiKeyAuth
func GetProductGroups(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.ProductGroup
	result := &[]models.ProductGroup{}

	models.SearchResourceQuery(&model, result, pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 更新主题产品
// @Produce  json
// @Tags 后台主题产品管理
// @Param id path int true "productGroup id"
// @Param params body ProductGroupParams true "productGroup params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/product_groups/{id} [put]
// @Security ApiKeyAuth
func UpdateProductGroup(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
		return
	}
	var err error
	var productGroupParams ProductGroupParams
	if err = c.ShouldBindJSON(&productGroupParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var productGroup models.ProductGroup
	id, _ := strconv.Atoi(c.Param("id"))
	productGroup.ID = id
	err = models.Find(&productGroup, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&productGroup, &productGroupParams)

	err = models.Save(&productGroup)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, productGroup)
}
