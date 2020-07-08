package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
)

type QueryLogisticParams struct {
	utils.Pagination
	ExpressCompany string `json:"q[express_company_cont]"`
	ExpressNo      string `json:"q[express_no_cont]"`
}

// @Summary 物流详情
// @Produce  json
// @Tags 后台物流管理
// @Param id path int true "logistic id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/logistics/{id} [get]
// @Security ApiKeyAuth
func GetLogistic(c *gin.Context) {
	var logistic models.Logistic
	id, _ := strconv.Atoi(c.Param("id"))
	logistic.ID = int(id)

	err := models.Find(&logistic, Query{Preloads: []string{"Order"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	apiHelpers.ResponseSuccess(c, logistic)
}

// @Summary 物流列表
// @Produce  json
// @Tags 后台物流管理
// @Param params query QueryLogisticParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/logistics [get]
// @Security ApiKeyAuth
func GetLogistics(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Logistic
	var result []models.Logistic

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q"), Preloads: []string{"Order"}}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)

}
