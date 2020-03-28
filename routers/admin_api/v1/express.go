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
)

type ExpressParams struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type QueryExpressParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
	Code string `json:"q[code_cont]"`
}

// @Summary 添加快递
// @Produce  json
// @Tags 后台快递管理
// @Param params body ExpressParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/expresses [post]
// @Security ApiKeyAuth
func AddExpress(c *gin.Context) {
	var err error
	var expressParams ExpressParams
	if err = c.ShouldBindJSON(&expressParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var express models.Express
	copier.Copy(&express, &expressParams)

	err = models.SaveResource(&express)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, express)
}

// @Summary 删除快递
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "express id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/expresses/{id} [delete]
// @Security ApiKeyAuth
func DeleteExpress(c *gin.Context) {
	var express models.Express
	id, _ := strconv.Atoi(c.Param("id"))
	express.ID = id

	err := models.DestroyResource(express, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 快递详情
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "express id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/expresses/{id} [get]
// @Security ApiKeyAuth
func GetExpress(c *gin.Context) {
	var express models.Express
	id, _ := strconv.Atoi(c.Param("id"))
	express.ID = id

	err := models.FindResource(&express, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	apiHelpers.ResponseSuccess(c, express)
}

// @Summary 快递列表
// @Produce  json
// @Tags 后台快递管理
// @Param params query QueryExpressParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/expresses [get]
// @Security ApiKeyAuth
func GetExpresses(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Express
	result := &[]models.Express{}

	models.SearchResourceQuery(&model, result, &pagination, c.QueryMap("q"))

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 更新快递
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "id"
// @Param params body ExpressParams true "express params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/expresses/{id} [put]
// @Security ApiKeyAuth
func UpdateExpress(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errors.New("id 不能为空"))
	}
	var err error
	var expressParams ExpressParams
	if err = c.ShouldBindJSON(&expressParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}

	var express models.Express

	id, _ := strconv.Atoi(c.Param("id"))
	express.ID = id
	err = models.SaveResource(&express)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	copier.Copy(&express, &expressParams)
	err = models.SaveResource(&express)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, express)
}
