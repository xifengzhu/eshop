package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	appApiHelper "github.com/xifengzhu/eshop/routers/app_api/api_helpers"
	"strconv"
)

// @Summary 领取优惠券
// @Produce  json
// @Tags 优惠券
// @Param code query string true "优惠券模板code"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/coupons/receive [post]
// @Security ApiKeyAuth
func CaptchCoupon(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	var template models.CouponTemplate

	parmMap := map[string]interface{}{"code": c.Query("code")}
	err := models.Where(Query{Conditions: parmMap}).Find(&template).Error
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	success, err := user.CatchCoupon(template.ID)
	if !success {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 获取用户的优惠券
// @Produce  json
// @Tags 优惠券
// @Param status query string true "优惠券状态" Enums(actived, expired, used) default(actived)
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/coupons [get]
// @Security ApiKeyAuth
func GetCoupons(c *gin.Context) {
	user := appApiHelper.CurrentUser(c)
	pagination := apiHelpers.SetDefaultPagination(c)
	var model models.Coupon
	result := &[]models.Coupon{}

	userIDStr := strconv.Itoa(user.ID)
	parmMap := map[string]string{"user_id": userIDStr, "state": c.DefaultQuery("status", "actived")}
	models.Search(&model, &Search{Pagination: pagination, Conditions: parmMap}, result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}
