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

// @Summary 添加分类
// @Produce  json
// @Tags 后台分类管理
// @Param params body params.CategoryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/categories [post]
// @Security ApiKeyAuth
func AddCategory(c *gin.Context) {
	var err error
	var categoryParams CategoryParams
	if err := ValidateParams(c, &categoryParams, "json"); err != nil {
		return
	}

	var category Category
	copier.Copy(&category, &categoryParams)

	err = Save(&category)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, category)
}

// @Summary 删除分类
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/categories/{id} [delete]
// @Security ApiKeyAuth
func DeleteCategory(c *gin.Context) {
	var category Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id

	var callbacks []func()
	callbacks = append(callbacks, category.RemoveChildrenRefer)
	err := DestroyWithCallbacks(&category, Options{Callbacks: callbacks})
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 分类详情
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/categories/{id} [get]
// @Security ApiKeyAuth
func GetCategory(c *gin.Context) {
	var category Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id
	err := Find(&category, Options{
		Preloads: []string{"Parent", "Children"},
	})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, category)
}

// @Summary 分类列表
// @Produce  json
// @Tags 后台分类管理
// @Param params query params.QueryCategoryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/categories [get]
// @Security ApiKeyAuth
func GetCategories(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model Category
	result := &[]Category{}

	Search(&model, &SearchParams{
		Pagination: pagination,
		Conditions: c.QueryMap("q"),
	}, &result)

	response := Collection{
		Pagination: pagination,
		List:       result,
	}

	ResponseSuccess(c, response)
}

// @Summary 更新分类
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Param params body params.CategoryParams true "category params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/categories/{id} [put]
// @Security ApiKeyAuth
func UpdateCategory(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
		return
	}
	var err error
	var categoryParams CategoryParams

	if err := ValidateParams(c, &categoryParams, "json"); err != nil {
		return
	}

	var category Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id
	err = Find(&category, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	copier.Copy(&category, &categoryParams)

	err = Save(&category)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, category)
}
