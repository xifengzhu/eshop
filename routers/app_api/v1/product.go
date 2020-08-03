package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"github.com/xifengzhu/eshop/routers/app_api/entities"
	"strconv"
)

type ProductQueryParams struct {
	utils.Pagination
	CategoryName string `form:"category_name"`
	Keyword      string `form:"keyword" validate:"required_without=CategoryName"`
}

// @Summary 获取产品列表
// @Produce  json
// @Tags 产品
// @Param params query ProductQueryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/products [get]
func GetProducts(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Product
	products := []models.Product{}

	var err error
	var productParams ProductQueryParams
	if err = apiHelpers.ValidateParams(c, &productParams, "query"); err != nil {
		return
	}

	categoryName := productParams.CategoryName
	if categoryName != "" {
		var category models.Category

		parmMap := map[string]interface{}{"name": categoryName}
		models.Where(Query{Conditions: parmMap}).Find(&category)
		products = category.GetCategoryProducts(pagination)
	} else {
		parmMap := map[string]string{"name_cont": productParams.Keyword}
		models.Search(&model, &Search{Pagination: pagination, Conditions: parmMap}, &products)
	}

	result := transferProductToEntity(products)
	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)

}

// @Summary 获取推荐列表
// @Produce  json
// @Tags 产品
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/recommend_products [get]
func GetRecommendProducts(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)
	products := &[]models.Product{}
	response := apiHelpers.Collection{Pagination: pagination, List: products}
	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 批量获取产品
// @Produce  json
// @Tags 产品
// @Param ids query []int true "产品ids"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/batch_products [get]
func BatchProducts(c *gin.Context) {
	parmSlice := c.QueryArray("ids")[:]
	var products []models.Product
	models.Where(Query{Conditions: parmSlice}).Find(&products)

	productEntities := transferProductToEntity(products)
	apiHelpers.ResponseSuccess(c, productEntities)
}

// @Summary 获取产品详情
// @Produce  json
// @Tags 产品
// @Param id path int true "id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/products/{id} [get]
func GetProduct(c *gin.Context) {
	var product models.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product.ID = id

	err := models.Find(&product, Query{Preloads: []string{"Goodses"}})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	product.SetSpecifications()
	var productDetail []entities.ProductDetailEntity
	copier.Copy(&productDetail, &product)

	apiHelpers.ResponseSuccess(c, productDetail[0])
}

func transferProductToEntity(products []models.Product) (productEntities []entities.ProductEntity) {
	for _, d_product := range products {
		var productEntity entities.ProductEntity
		copier.Copy(&productEntity, &d_product)
		productEntities = append(productEntities, productEntity)
	}
	return
}
