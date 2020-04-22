package v1

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"io/ioutil"
)

type AppSettingParams struct {
	WxappId       string `json:"app_id" binding:"required"`
	AppName       string `json:"app_name" binding:"required"`
	AppSecret     string `json:"app_secret" binding:"required"`
	ServicePhone  string `json:"sevice_phone"`
	Mchid         string `json:"mchid" binding:"required"`
	Apikey        string `json:"api_key" binding:"required"`
	ApiClientCert string `json:"api_client_cert,omitempty"`
	NotifyUrl     string `json:"notify_url,omitempty"`
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

// @Summary 更新证书
// @Produce  json
// @Tags 后台配置管理
// @Accept  multipart/form-data
// @Param api_client_cert formData file true "wechat pay certification"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/app_setting/cert [post]
// @Security ApiKeyAuth
func UpdateWechatCert(c *gin.Context) {
	var err error
	var setting models.AppSetting
	setting.Current()

	// read file base64
	fh, _ := c.FormFile("api_client_cert")
	file, _ := fh.Open()
	defer file.Close()
	bytes, _ := ioutil.ReadAll(file)
	base64Str := base64.StdEncoding.EncodeToString(bytes)

	setting.ApiClientCert = base64Str
	err = models.SaveResource(&setting)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
		return
	}
	apiHelpers.ResponseOK(c)
}
