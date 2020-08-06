package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/app_api/helpers"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 领取优惠券
// @Produce  json
// @Tags 优惠券
// @Param code query string true "优惠券模板code"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/coupons/receive [post]
// @Security ApiKeyAuth
func CaptchCoupon(c *gin.Context) {
	user := CurrentUser(c)
	var template CouponTemplate

	parmMap := map[string]interface{}{"code": c.Query("code")}
	err := Where(Options{Conditions: parmMap}).Find(&template).Error
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	success, err := user.CatchCoupon(template.ID)
	if !success {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	ResponseOK(c)
}

// @Summary 获取用户的优惠券
// @Produce  json
// @Tags 优惠券
// @Param status query string true "优惠券状态" Enums(actived, expired, used) default(actived)
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/coupons [get]
// @Security ApiKeyAuth
func GetCoupons(c *gin.Context) {
	user := CurrentUser(c)
	pagination := SetDefaultPagination(c)
	var model Coupon
	result := &[]Coupon{}

	userIDStr := strconv.Itoa(user.ID)
	parmMap := map[string]string{"user_id": userIDStr, "state": c.DefaultQuery("status", "actived")}
	Search(&model, &SearchParams{Pagination: pagination, Conditions: parmMap}, result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}
