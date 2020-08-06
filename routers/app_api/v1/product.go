package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/app_api/present"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 获取产品列表
// @Produce  json
// @Tags 产品
// @Param params query params.ProductQueryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/products [get]
func GetProducts(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model Product
	products := []Product{}

	var err error
	var productParams ProductQueryParams
	if err = ValidateParams(c, &productParams, "query"); err != nil {
		return
	}

	categoryName := productParams.CategoryName
	if categoryName != "" {
		var category Category

		parmMap := map[string]interface{}{"name": categoryName}
		Where(Options{Conditions: parmMap}).Find(&category)
		products = category.GetCategoryProducts(pagination)
	} else {
		parmMap := map[string]string{"name_cont": productParams.Keyword}
		Search(&model, &SearchParams{Pagination: pagination, Conditions: parmMap}, &products)
	}

	result := transferProductToEntity(products)
	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)

}

// @Summary 获取推荐列表
// @Produce  json
// @Tags 产品
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/recommend_products [get]
func GetRecommendProducts(c *gin.Context) {
	pagination := SetDefaultPagination(c)
	products := &[]Product{}
	response := Collection{Pagination: pagination, List: products}
	ResponseSuccess(c, response)
}

// @Summary 批量获取产品
// @Produce  json
// @Tags 产品
// @Param ids query []int true "产品ids"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/batch_products [get]
func BatchProducts(c *gin.Context) {
	parmSlice := c.QueryArray("ids")[:]
	var products []Product
	Where(Options{Conditions: parmSlice}).Find(&products)

	productEntities := transferProductToEntity(products)
	ResponseSuccess(c, productEntities)
}

// @Summary 获取产品详情
// @Produce  json
// @Tags 产品
// @Param id path int true "id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/products/{id} [get]
func GetProduct(c *gin.Context) {
	var product Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := Find(&product, Options{Preloads: []string{"Goodses"}})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	product.SetSpecifications()
	var productDetail []ProductDetailEntity
	copier.Copy(&productDetail, &product)

	ResponseSuccess(c, productDetail[0])
}

func transferProductToEntity(products []Product) (productEntities []ProductEntity) {
	for _, d_product := range products {
		var productEntity ProductEntity
		copier.Copy(&productEntity, &d_product)
		productEntities = append(productEntities, productEntity)
	}
	return
}
