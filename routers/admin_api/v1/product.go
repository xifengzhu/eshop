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
	// "time"
)

type GoodsParams struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name"`
	Properties string  `json:"properties"`
	Images     string  `json:"images"`
	SkuNo      string  `json:"sku_no"`
	StockNum   int     `json:"stock_num"`
	Position   int     `json:"position"`
	Price      float32 `json:"price"`
	LinePrice  float32 `json:"line_price"`
	Weight     float32 `json:"weight"`
	Destroy    bool    `json:"_destroy,omitempty"`
}

type ProductParams struct {
	Name            string        `json:"name"`
	Content         string        `json:"content"`
	DeductStockType int           `json:"deduct_stock_type"`
	SalesInitial    int           `json:"sales_initial"`
	Position        int           `json:"position"`
	Price           float32       `json:"price"`
	IsOnline        bool          `json:"is_online"`
	DeliveryID      int           `json:"product_id"`
	Goodses         []GoodsParams `json:"goodses"`
}

type QueryProductParams struct {
	utils.Pagination
	IsOnline   bool    `json:"q[is_online_eq]"`
	Name       string  `json:"q[name_cont]"`
	Price_gteq float32 `json:"q[price_gteq]"`
	Price_lteq float32 `json:"q[price_lteq]"`
}

// @Summary 添加产品
// @Produce  json
// @Tags 后台产品管理
// @Param params body ProductParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products [post]
// @Security ApiKeyAuth
func AddProduct(c *gin.Context) {
	var err error
	var productParams ProductParams
	if err = c.ShouldBindJSON(&productParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var product models.Product
	copier.Copy(&product, &productParams)

	err = models.SaveResource(&product)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, product)
}

// @Summary 删除产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [delete]
// @Security ApiKeyAuth
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := models.DestroyResource(product, Query{Callbacks: []func(){product.RemoveGoodses}})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 产品列表
// @Produce  json
// @Tags 后台产品管理
// @Param params query QueryProductParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products [get]
// @Security ApiKeyAuth
func GetProductes(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Product
	result := &[]models.Product{}

	models.SearchResourceQuery(&model, result, pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 产品详情
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "product id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [get]
// @Security ApiKeyAuth
func GetProduct(c *gin.Context) {
	var product models.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := models.FindResource(&product, Query{Preloads: []string{"Goodses"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, product)
}

// @Summary 更新产品
// @Produce  json
// @Tags 后台产品管理
// @Param id path int true "id"
// @Param params body ProductParams true "product params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/products/{id} [put]
// @Security ApiKeyAuth
func UpdateProduct(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
	}
	var err error
	var productParams ProductParams
	if err = c.ShouldBindJSON(&productParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var product models.Product

	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err = models.FindResource(&product, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&product, &productParams)
	// reset product goodses
	product.Goodses = nil
	// recover the product id
	product.ID = id
	copier.Copy(&product.Goodses, &productParams.Goodses)

	err = product.NestUpdate()
	models.FindResource(&product, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, product)
}
