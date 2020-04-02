package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
)

type AddressParams struct {
	UserID     int    `json:"user_id" binding:"required"`
	RegionID   int    `json:"region_id" binding:"required"`
	ProvinceID int    `json:"province_id" binding:"required"`
	CityID     int    `json:"city_id" binding:"required"`
	Detail     string `json:"detail" binding:"required"`
	isDefault  bool   `json:"is_default" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
}

// @Summary 新增地址
// @Produce  json
// @Tags 地址
// @Param params body AddressParams true "address params"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses [post]
// @Security ApiKeyAuth
func AddAddress(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	var addressParams AddressParams
	addressParams.UserID = user.ID
	if err := c.ShouldBindJSON(&addressParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}

	var address models.Address
	copier.Copy(&address, &addressParams)

	err := models.SaveResource(&address)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, address)
}

// @Summary 编辑地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Param params body AddressParams true "address attributes"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses/{id} [put]
// @Security ApiKeyAuth
func EditAddress(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	addressID, _ := strconv.Atoi(c.Param("id"))
	var address models.Address
	address.ID = addressID
	if address.Exist() == false {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, errors.New("资源不存在"))
		return
	}

	var addressParams AddressParams
	addressParams.UserID = user.ID
	if err := c.ShouldBindJSON(&addressParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}

	copier.Copy(&address, &addressParams)
	address.ID = addressID

	err := models.SaveResource(&address)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}
	apiHelpers.ResponseSuccess(c, address)
}

// @Summary 获取用户的收获地址
// @Produce  json
// @Tags 地址
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses [get]
// @Security ApiKeyAuth
func GetAddresses(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	var addresses []models.Address
	parmMap := map[string]interface{}{"user_id": user.ID}
	models.WhereResources(&addresses, Query{Conditions: parmMap})
	apiHelpers.ResponseSuccess(c, addresses)
}

// @Summary 通过id获取地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses/{id} [get]
// @Security ApiKeyAuth
func GetAddress(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	ID, _ := strconv.Atoi(c.Param("id"))

	var addresses []models.Address
	parmMap := map[string]interface{}{"id": ID, "user_id": user.ID}
	err := models.WhereResources(&addresses, Query{Conditions: parmMap})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
	} else {
		apiHelpers.ResponseSuccess(c, addresses[0])
	}
}

// @Summary 删除地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses/{id} [delete]
// @Security ApiKeyAuth
func DeleteAddress(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)

	ID, _ := strconv.Atoi(c.Param("id"))
	address, err := user.GetAddressByID(ID)
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
		return
	}

	models.DestroyResource(&address, Query{})
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 获取所有的省份信息
// @Produce  json
// @Tags 地址
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/provinces [get]
func GetProvinces(c *gin.Context) {
	var provinces []models.Province
	models.AllResource(&provinces, Query{})
	response := apiHelpers.Collection{List: provinces}
	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 根据省ID获取城市列表信息
// @Produce  json
// @Tags 地址
// @Param province_id query string true "province id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/cities [get]
func GetCities(c *gin.Context) {
	provinceID := c.Query("province_id")
	checkEmptyParams(provinceID, c)
	var cities []models.City
	parmMap := map[string]interface{}{"province_id": provinceID}
	models.WhereResources(&cities, Query{Conditions: parmMap})
	apiHelpers.ResponseSuccess(c, cities)
}

// @Summary 根据市ID获取下级城市列表信息
// @Produce  json
// @Tags 地址
// @Param city_id query string true "city id"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/regions [get]
func GetRegions(c *gin.Context) {
	cityID := c.Query("city_id")
	checkEmptyParams(cityID, c)
	var regions []models.Region
	parmMap := map[string]interface{}{"city_id": cityID}
	models.WhereResources(&regions, Query{Conditions: parmMap})

	apiHelpers.ResponseSuccess(c, regions)
}

func checkEmptyParams(param string, c *gin.Context) {
	if param == "" {
		response := apiHelpers.Response{Code: e.INVALID_PARAMS, Data: nil, Msg: e.GetMsg(e.INVALID_PARAMS)}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}
