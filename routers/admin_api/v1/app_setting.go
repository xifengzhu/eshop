package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	// "github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	// "strconv"
)

type AppSettingParams struct {
	AppName      string `json:"app_name" binding:"required"`
	AppSecret    string `json:"app_secret" binding:"required"`
	ServicePhone string `json:"sevice_phone" binding:"required"`
	Mchid        string `json:"mchid" binding:"required"`
	Apikey       string `json:"api_key" binding:"required"`
	NotifyUrl    string `json:"notify_url" binding:"required"`
}

// @Summary 配置详情
// @Produce  json
// @Tags 后台配置管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/app_setting [get]
// @Security ApiKeyAuth
func GetAppSetting(c *gin.Context) {
	var setting models.AppSetting

	setting.Current()

	apiHelpers.ResponseSuccess(c, setting)
}

// @Summary 更新配置
// @Produce  json
// @Tags 后台配置管理
// @Param params body AppSettingParams true "app_setting params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/app_setting [put]
// @Security ApiKeyAuth
func UpdateAppSetting(c *gin.Context) {
	var err error
	var settingParams AppSettingParams
	if err = c.ShouldBindJSON(&settingParams); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}

	var setting models.AppSetting

	copier.Copy(&setting, &settingParams)

	err = setting.CreateOrUpdate()
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseSuccess(c, setting)
}
