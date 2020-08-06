package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"log"
	"strconv"
)

// @Summary 添加运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param params body params.DeliveryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/deliveries [post]
// @Security ApiKeyAuth
func AddDelivery(c *gin.Context) {
	var err error
	var deliveryParams DeliveryParams
	if err := ValidateParams(c, &deliveryParams, "json"); err != nil {
		return
	}

	var delivery Delivery
	copier.Copy(&delivery, &deliveryParams)

	err = Save(&delivery)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, delivery)
}

// @Summary 删除运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "delivery id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/deliveries/{id} [delete]
// @Security ApiKeyAuth
func DeleteDelivery(c *gin.Context) {
	var delivery Delivery
	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = id

	var callbacks []func()
	callbacks = append(callbacks, delivery.DestroyRules)
	err := DestroyWithCallbacks(&delivery, Options{Callbacks: callbacks})

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 运费模板详情
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "delivery id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/deliveries/{id} [get]
// @Security ApiKeyAuth
func GetDelivery(c *gin.Context) {
	var delivery Delivery
	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = int(id)

	err := Find(&delivery, Options{Preloads: []string{"DeliveryRules"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, delivery)
}

// @Summary 运费模板列表
// @Produce  json
// @Tags 后台运费模板管理
// @Param params query params.QueryDeliveryParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/deliveries [get]
// @Security ApiKeyAuth
func GetDeliveries(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model Delivery
	result := &[]Delivery{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)

}

// @Summary 更新运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "id"
// @Param params body params.DeliveryParams true "delivery params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/deliveries/{id} [put]
// @Security ApiKeyAuth
func UpdateDelivery(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var deliveryParams DeliveryParams
	if err := ValidateParams(c, &deliveryParams, "json"); err != nil {
		return
	}

	var delivery Delivery

	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = id

	err = Find(&delivery, Options{})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	log.Println("=======deliveryParams====", deliveryParams)
	copier.Copy(&delivery, &deliveryParams)
	// reset delivery rules
	delivery.DeliveryRules = nil
	// recover the delivery id
	delivery.ID = id
	copier.Copy(&delivery.DeliveryRules, &deliveryParams.DeliveryRules)
	log.Println("=======delivery====", delivery)
	err = delivery.NestUpdate()
	delivery.Reload()
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, delivery)
}
