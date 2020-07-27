package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"log"
	"strconv"
)

type DeliveryRuleParams struct {
	ID            int         `json:"id,omitempty"`
	First         float64     `json:"first" ` // 首件/首重
	FirstFee      float64     `json:"first_fee" `
	Additional    float64     `json:"additional" `     // 续件/续重
	AdditionalFee float64     `json:"additional_fee" ` // 续件/续重
	Region        models.JSON `json:"region" `         // 可配送区域(省id集)
	Destroy       bool        `json:"_destroy,omitempty"`
	Position      int         `json:"position"`
}

type DeliveryParams struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty" validate:"required"`
	Way           int                  `json:"way" validate:"required"` // 1 为按件计费 2 按重量计费
	DeliveryRules []DeliveryRuleParams `json:"delivery_rules" validate:"required,dive"`
}

type QueryDeliveryParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}

// @Summary 添加运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param params body DeliveryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/deliveries [post]
// @Security ApiKeyAuth
func AddDelivery(c *gin.Context) {
	var err error
	var deliveryParams DeliveryParams
	if err := apiHelpers.ValidateParams(c, &deliveryParams); err != nil {
		return
	}

	var delivery models.Delivery
	copier.Copy(&delivery, &deliveryParams)

	err = models.Save(&delivery)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, delivery)
}

// @Summary 删除运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "delivery id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/deliveries/{id} [delete]
// @Security ApiKeyAuth
func DeleteDelivery(c *gin.Context) {
	var delivery models.Delivery
	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = id

	var callbacks []func()
	callbacks = append(callbacks, delivery.DestroyRules)
	err := models.DestroyWithCallbacks(&delivery, Query{Callbacks: callbacks})

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 运费模板详情
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "delivery id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/deliveries/{id} [get]
// @Security ApiKeyAuth
func GetDelivery(c *gin.Context) {
	var delivery models.Delivery
	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = int(id)

	err := models.Find(&delivery, Query{Preloads: []string{"DeliveryRules"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	apiHelpers.ResponseSuccess(c, delivery)
}

// @Summary 运费模板列表
// @Produce  json
// @Tags 后台运费模板管理
// @Param params query QueryDeliveryParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/deliveries [get]
// @Security ApiKeyAuth
func GetDeliveries(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.Delivery
	result := &[]models.Delivery{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)

}

// @Summary 更新运费模板
// @Produce  json
// @Tags 后台运费模板管理
// @Param id path int true "id"
// @Param params body DeliveryParams true "delivery params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/deliveries/{id} [put]
// @Security ApiKeyAuth
func UpdateDelivery(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var deliveryParams DeliveryParams
	if err := apiHelpers.ValidateParams(c, &deliveryParams); err != nil {
		return
	}

	var delivery models.Delivery

	id, _ := strconv.Atoi(c.Param("id"))
	delivery.ID = id

	err = models.Find(&delivery, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, delivery)
}
