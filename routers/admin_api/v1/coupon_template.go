package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
	"time"
)

type QueryCouponTemplateParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}

type Config struct {
	MinAmount    float64    `json:"min_amount" validate:"required,gte=0"`
	ResourceType string     `json:"resource_type,omitempty"`
	Resources    []int      `json:"resources,omitempty" validate:"required_with=ResourceType"`
	DateType     string     `json:"date_type" validate:"required,oneof='fix_term' 'time_range'"`
	FixTerm      int        `json:"fix_term,omitempty" validate:"rfe=DateType:fix_term"`
	StartAt      *time.Time `json:"start_at,omitempty" validate:"required_with=EndAt,rfe=DateType:time_range"`
	EndAt        *time.Time `json:"end_at,omitempty" validate:"required_with=StartAt,rfe=DateType:time_range""`

	ReduceAmount float64 `json:"reduce_amount,omitempty" validate:"required_without=Percentage"`
	Percentage   int     `json:"percentage,omitempty" validate:"required_without=ReduceAmount"`
}

type CouponTemplateParams struct {
	Code       string     `json:"code"`
	Name       string     `json:"name" validate:"required"`
	Kind       string     `json:"kind" validate:"required,oneof='fixed_amount' 'percentage'"`
	Creator    string     `json:"creator"`
	Stock      int        `json:"stock" validate:"required"`
	CatchLimit int        `json:"catch_limit" validate:"required"`
	StartAt    *time.Time `json:"start_at" validate:"required"`
	EndAt      *time.Time `json:"end_at" validate:"required,gtfield=StartAt"`
	Configs    Config     `json:"configs" validate:"required,dive"`
}

// @Summary 添加优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param params body CouponTemplateParams true "coupon template params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/coupon_templates [post]
// @Security ApiKeyAuth
func AddCouponTemplate(c *gin.Context) {
	var err error
	var templateParams CouponTemplateParams
	if err := apiHelpers.ValidateParams(c, &templateParams, "json"); err != nil {
		return
	}

	var template models.CouponTemplate
	copier.Copy(&template, &templateParams)
	jsonConfigs, _ := json.Marshal(templateParams.Configs)
	template.Configs = []byte(jsonConfigs)
	err = models.Save(&template)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, template)
}

// @Summary 删除优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "template id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [delete]
// @Security ApiKeyAuth
func DeleteCouponTemplate(c *gin.Context) {
	var template models.CouponTemplate
	id, _ := strconv.Atoi(c.Param("id"))
	template.ID = id

	err := models.DestroyWithCallbacks(&template, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, nil)
}

// @Summary 优惠券模板详情
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "template id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [get]
// @Security ApiKeyAuth
func GetCouponTemplate(c *gin.Context) {
	var template models.CouponTemplate
	id, _ := strconv.Atoi(c.Param("id"))
	template.ID = int(id)

	err := models.Find(&template, Query{Preloads: []string{"Coupons"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	apiHelpers.ResponseSuccess(c, template)
}

// @Summary 优惠券模板列表
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param params query QueryCouponTemplateParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/coupon_templates [get]
// @Security ApiKeyAuth
func GetCouponTemplates(c *gin.Context) {

	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.CouponTemplate
	result := &[]models.CouponTemplate{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)

}

// @Summary 更新优惠券模板
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "id"
// @Param params body CouponTemplateParams true "template params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/coupon_templates/{id} [put]
// @Security ApiKeyAuth
func UpdateCouponTemplate(c *gin.Context) {
	if c.Param("id") == "" {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "id 不能为空")
	}
	var err error
	var params CouponTemplateParams
	if err = apiHelpers.ValidateParams(c, &params, "json"); err != nil {
		return
	}

	var template models.CouponTemplate
	template.ID, _ = strconv.Atoi(c.Param("id"))
	err = models.Find(&template, Query{})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	changedAttrs := models.CouponTemplate{}
	copier.Copy(&changedAttrs, &params)
	err = models.Update(&template, &changedAttrs)

	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, template)
}

// @Summary 创建优惠券
// @Produce  json
// @Tags 后台优惠券模板管理
// @Param id path int true "id"
// @Param qty query int true "quantity"
// @Success 200 {object} apiHelpers.Response
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
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, errMsg)
	}

	var template models.CouponTemplate
	template.ID, _ = strconv.Atoi(c.Param("id"))
	err := models.Find(&template, Query{})

	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
	}

	qty, _ := strconv.Atoi(c.Query("qty"))
	template.GenerateCoupon(qty)

	apiHelpers.ResponseOK(c)
}
