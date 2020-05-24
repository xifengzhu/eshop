package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
	// "log"
	"net/http"
	"strconv"
)

type AddressParams struct {
	UserID    int    `json:"user_id"`
	Region    string `json:"region" binding:"required"`
	Province  string `json:"province" binding:"required"`
	City      string `json:"city" binding:"required"`
	Detail    string `json:"detail" binding:"required"`
	isDefault bool   `json:"is_default" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Receiver  string `json:"receiver" binding:"required"`
}

type AddressQueryParams struct {
	utils.Pagination
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
		return
	}

	var address models.Address
	copier.Copy(&address, &addressParams)

	err := models.Save(&address)
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

	changedAttrs := models.Address{}
	copier.Copy(&changedAttrs, &addressParams)

	err := models.Update(&address, changedAttrs)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}
	apiHelpers.ResponseSuccess(c, address)
}

// @Summary 获取用户的收获地址
// @Produce  json
// @Tags 地址
// @Param params body AddressQueryParams true "address pagination"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/addresses [get]
// @Security ApiKeyAuth
func GetAddresses(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	pagination := apiHelpers.SetDefaultPagination(c)
	var model models.Address
	result := &[]models.Address{}

	userIDStr := strconv.Itoa(user.ID)
	parmMap := map[string]string{"user_id": userIDStr}
	models.Search(&model, &Search{Pagination: pagination, Conditions: parmMap}, result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
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

	addresses, err := user.GetAddressByID(ID)

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err)
	} else {
		apiHelpers.ResponseSuccess(c, addresses)
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

	models.Destroy(&address)
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 获取所有的省份信息
// @Produce  json
// @Tags 地址
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/provinces [get]
func GetProvinces(c *gin.Context) {
	var provinces []models.Province
	models.All(&provinces, Query{})
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
	models.Where(Query{Conditions: parmMap}).Find(&cities)
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
	models.Where(Query{Conditions: parmMap}).Find(&regions)

	apiHelpers.ResponseSuccess(c, regions)
}

func checkEmptyParams(param string, c *gin.Context) {
	if param == "" {
		response := apiHelpers.Response{Code: e.INVALID_PARAMS, Data: nil, Msg: e.GetMsg(e.INVALID_PARAMS)}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
}
