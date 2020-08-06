package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	_ "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 物流详情
// @Produce  json
// @Tags 后台物流管理
// @Param id path int true "logistic id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/logistics/{id} [get]
// @Security ApiKeyAuth
func GetLogistic(c *gin.Context) {
	var logistic Logistic
	id, _ := strconv.Atoi(c.Param("id"))
	logistic.ID = int(id)

	err := Find(&logistic, Options{Preloads: []string{"Order"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, logistic)
}

// @Summary 物流列表
// @Produce  json
// @Tags 后台物流管理
// @Param params query params.QueryLogisticParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/logistics [get]
// @Security ApiKeyAuth
func GetLogistics(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model Logistic
	var result []Logistic

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q"), Preloads: []string{"Order"}}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)

}
