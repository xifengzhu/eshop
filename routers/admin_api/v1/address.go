package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"net/http"
)

// @Summary 获取所有的省份信息
// @Produce  json
// @Tags 后台管理地址
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/provinces [get]
func GetProvinces(c *gin.Context) {
	var provinces []Province
	All(&provinces, Options{})
	response := Collection{List: provinces}
	ResponseSuccess(c, response)
}

// @Summary 获取用户收货地址
// @Produce  json
// @Tags 后台管理地址
// @Param user_id query int true "user id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/addresses [get]
func GetUserAdddresses(c *gin.Context) {
	userID := c.Query("user_id")
	var addresses []Address
	parmMap := map[string]interface{}{"user_id": userID}
	Where(Options{Conditions: parmMap}).Find(&addresses)
	ResponseSuccess(c, addresses)
}

// @Summary 根据省ID获取城市列表信息
// @Produce  json
// @Tags 后台管理地址
// @Param province_id query string true "province id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/cities [get]
func GetCities(c *gin.Context) {
	provinceID := c.Query("province_id")
	checkEmptyParams(provinceID, c)
	var cities []City
	parmMap := map[string]interface{}{"province_id": provinceID}
	Where(Options{Conditions: parmMap}).Find(&cities)
	ResponseSuccess(c, cities)
}

// @Summary 根据市ID获取下级城市列表信息
// @Produce  json
// @Tags 后台管理地址
// @Param city_id query string true "city id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/regions [get]
func GetRegions(c *gin.Context) {
	cityID := c.Query("city_id")
	checkEmptyParams(cityID, c)
	var regions []Region
	parmMap := map[string]interface{}{"city_id": cityID}
	Where(Options{Conditions: parmMap}).Find(&regions)

	ResponseSuccess(c, regions)
}

func checkEmptyParams(param string, c *gin.Context) {
	if param == "" {
		response := Response{Code: e.INVALID_PARAMS, Data: nil, Msg: e.GetMsg(e.INVALID_PARAMS)}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}
