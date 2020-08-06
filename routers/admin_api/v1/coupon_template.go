package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 添加优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param params body params.CouponTemplateParams true "coupon template params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates [post]
// @Security ApiKeyAuth
func AddCouponTemplate(c *gin.Context) {
	var err error
	var templateParams CouponTemplateParams
	if err := ValidateParams(c, &templateParams, "json"); err != nil {
		return
	}

	var template CouponTemplate
	copier.Copy(&template, &templateParams)
	jsonConfigs, _ := json.Marshal(templateParams.Configs)
	template.Configs = []byte(jsonConfigs)
	err = Save(&template)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, template)
}

// @Summary 删除优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "template id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [delete]
// @Security ApiKeyAuth
func DeleteCouponTemplate(c *gin.Context) {
	var template CouponTemplate
	id, _ := strconv.Atoi(c.Param("id"))
	template.ID = id

	err := DestroyWithCallbacks(&template, Options{})

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, nil)
}

// @Summary 优惠券模板详情
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "template id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [get]
// @Security ApiKeyAuth
func GetCouponTemplate(c *gin.Context) {
	var template CouponTemplate
	id, _ := strconv.Atoi(c.Param("id"))
	template.ID = int(id)

	err := Find(&template, Options{Preloads: []string{"Coupons"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, template)
}

// @Summary 优惠券模板列表
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param params query params.QueryCouponTemplateParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates [get]
// @Security ApiKeyAuth
func GetCouponTemplates(c *gin.Context) {

	pagination := SetDefaultPagination(c)

	var model CouponTemplate
	result := &[]CouponTemplate{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)

}

// @Summary 更新优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "id"
// @Param params body params.CouponTemplateParams true "template params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [put]
// @Security ApiKeyAuth
func UpdateCouponTemplate(c *gin.Context) {
	if c.Param("id") == "" {
		ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var params CouponTemplateParams
	if err = ValidateParams(c, &params, "json"); err != nil {
		return
	}

	var template CouponTemplate
	template.ID, _ = strconv.Atoi(c.Param("id"))
	err = Find(&template, Options{})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	changedAttrs := CouponTemplate{}
	copier.Copy(&changedAttrs, &params)
	err = Update(&template, &changedAttrs)

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseSuccess(c, template)
}

// @Summary 创建优惠券
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "id"
// @Param qty query int true "quantity"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/coupon_templates/{id}/generate_coupons [post]
// @Security ApiKeyAuth
func GenerateCoupons(c *gin.Context) {
	var errMsg string
	if c.Param("id") == "" {
		errMsg = "id 不能为空"
	}

	if c.Query("qty") == "" {
		errMsg = "数量不能为空"
	}

	if errMsg != "" {
		ResponseError(c, e.INVALID_PARAMS, errMsg)
	}

	var template CouponTemplate
	template.ID, _ = strconv.Atoi(c.Param("id"))
	err := Find(&template, Options{})

	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
	}

	qty, _ := strconv.Atoi(c.Query("qty"))
	template.GenerateCoupon(qty)

	ResponseOK(c)
}
