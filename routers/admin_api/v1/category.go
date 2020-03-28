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

type CategoryParams struct {
	Name     string `binding:"required" json:"name"`
	Position int    `json:"position"`
	ParentID int    `json:"parent_id"`
	Image    string `json:"image"`
}

type QueryCategoryParams struct {
	utils.Pagination
	Name            string    `uri:"q[name]"`
	Created_at_gteq time.Time `uri:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `uri:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
}

// @Summary 添加分类
// @Produce  json
// @Tags 后台分类管理
// @Param params body CategoryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/categories [post]
// @Security ApiKeyAuth
func AddCategory(c *gin.Context) {
	var err error
	var categoryParams CategoryParams
	if err = c.ShouldBind(&categoryParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var category models.Category
	copier.Copy(&category, &categoryParams)

	err = models.SaveResource(&category)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, category)
}

// @Summary 删除分类
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/categories/{id} [delete]
// @Security ApiKeyAuth
func DeleteCategory(c *gin.Context) {
	var category models.Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id

	var callbacks []func()
	callbacks = append(callbacks, category.RemoveChildrenRefer)
	err := models.DestroyResource(&category, Query{Callbacks: callbacks})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 分类详情
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/categories/{id} [get]
// @Security ApiKeyAuth
func GetCategory(c *gin.Context) {
	var category models.Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id
	err := models.FindResource(&category, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, category)
}

// @Summary 分类列表
// @Produce  json
// @Tags 后台分类管理
// @Param params query QueryCategoryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/categories [get]
// @Security ApiKeyAuth
func GetCategories(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Category
	result := &[]models.Category{}

	models.SearchResourceQuery(&model, result, &pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 更新分类
// @Produce  json
// @Tags 后台分类管理
// @Param id path int true "category id"
// @Param params body CategoryParams true "category params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/categories/{id} [put]
// @Security ApiKeyAuth
func UpdateCategory(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
		return
	}
	var err error
	var categoryParams CategoryParams
	if err = c.ShouldBindJSON(&categoryParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var category models.Category
	id, _ := strconv.Atoi(c.Param("id"))
	category.ID = id
	err = models.FindResource(&category, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&category, &categoryParams)

	err = models.SaveResource(&category)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, category)
}
