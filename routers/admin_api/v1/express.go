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

// @Summary 添加快递
// @Produce  json
// @Tags 后台快递管理
// @Param params body params.ExpressParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/expresses [post]
// @Security ApiKeyAuth
func AddExpress(c *gin.Context) {
	var err error
	var expressParams ExpressParams
	if err := ValidateParams(c, &expressParams, "json"); err != nil {
		return
	}

	var express Express
	copier.Copy(&express, &expressParams)

	err = Save(&express)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, express)
}

// @Summary 删除快递
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "express id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/expresses/{id} [delete]
// @Security ApiKeyAuth
func DeleteExpress(c *gin.Context) {
	var express Express
	id, _ := strconv.Atoi(c.Param("id"))
	express.ID = id

	err := Destroy(express)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 快递详情
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "express id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/expresses/{id} [get]
// @Security ApiKeyAuth
func GetExpress(c *gin.Context) {
	var express Express
	id, _ := strconv.Atoi(c.Param("id"))
	express.ID = id

	err := Find(&express, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, express)
}

// @Summary 快递列表
// @Produce  json
// @Tags 后台快递管理
// @Param params query params.QueryExpressParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/expresses [get]
// @Security ApiKeyAuth
func GetExpresses(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model Express
	result := &[]Express{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 更新快递
// @Produce  json
// @Tags 后台快递管理
// @Param id path int true "id"
// @Param params body params.ExpressParams true "express params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/expresses/{id} [put]
// @Security ApiKeyAuth
func UpdateExpress(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var expressParams ExpressParams
	if err = ValidateParams(c, &expressParams, "json"); err != nil {
		return
	}

	var express Express
	express.ID, _ = strconv.Atoi(c.Param("id"))
	err = Find(&express, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	changedAttrs := Express{}
	copier.Copy(&changedAttrs, &expressParams)
	err = Update(&express, &changedAttrs)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, express)
}
