package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/app_api/helpers"
	. "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 新增地址
// @Produce  json
// @Tags 地址
// @Param params body params.AddressParams true "address params"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/addresses [post]
// @Security ApiKeyAuth
func AddAddress(c *gin.Context) {
	user := CurrentUser(c)
	var params AddressParams
	params.UserID = user.ID

	if err := ValidateParams(c, &params, "json"); err != nil {
		return
	}

	var address Address
	copier.Copy(&address, &params)

	err := Save(&address)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, address)
}

// @Summary 编辑地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Param params body params.AddressParams true "address attributes"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/addresses/{id} [put]
// @Security ApiKeyAuth
func EditAddress(c *gin.Context) {
	user := CurrentUser(c)

	addressID, _ := strconv.Atoi(c.Param("id"))
	var address Address
	address.ID = addressID
	if address.Exist() == false {
		ResponseError(c, e.ERROR_NOT_EXIST, "资源不存在")
		return
	}

	var params AddressParams
	params.UserID = user.ID

	if err := ValidateParams(c, &params, "json"); err != nil {
		return
	}

	changedAttrs := Address{}
	copier.Copy(&changedAttrs, &params)

	err := Update(&address, changedAttrs)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
	}
	ResponseSuccess(c, address)
}

// @Summary 获取用户的收获地址
// @Produce  json
// @Tags 地址
// @Param params query params.AddressQueryParams true "address pagination"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/addresses [get]
// @Security ApiKeyAuth
func GetAddresses(c *gin.Context) {
	user := CurrentUser(c)
	pagination := SetDefaultPagination(c)
	var model Address
	result := &[]Address{}

	userIDStr := strconv.Itoa(user.ID)
	parmMap := map[string]string{"user_id": userIDStr}
	Search(&model, &SearchParams{Pagination: pagination, Conditions: parmMap}, result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 通过id获取地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/addresses/{id} [get]
// @Security ApiKeyAuth
func GetAddress(c *gin.Context) {
	user := CurrentUser(c)
	ID, _ := strconv.Atoi(c.Param("id"))

	addresses, err := user.GetAddressByID(ID)

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
	} else {
		ResponseSuccess(c, addresses)
	}
}

// @Summary 删除地址
// @Produce  json
// @Tags 地址
// @Param id path int true "id"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/addresses/{id} [delete]
// @Security ApiKeyAuth
func DeleteAddress(c *gin.Context) {
	user := CurrentUser(c)

	ID, _ := strconv.Atoi(c.Param("id"))
	address, err := user.GetAddressByID(ID)
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	Destroy(&address)
	ResponseSuccess(c, nil)
}

// @Summary 获取所有的省份信息
// @Produce  json
// @Tags 地址
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/provinces [get]
func GetProvinces(c *gin.Context) {
	var provinces []Province
	All(&provinces, Options{Preloads: []string{"Cities", "Cities.Regions"}})
	response := Collection{List: provinces}
	ResponseSuccess(c, response)
}
