package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	"github.com/xifengzhu/eshop/routers/admin_api/entities"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"io"
	"os"
	"path/filepath"
)

// @Summary 微信支付配置详情
// @Produce  json
// @Tags 后台配置管理
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxpay_setting [get]
// @Security ApiKeyAuth
func GetWxpaySetting(c *gin.Context) {
	var setting models.WxpaySetting

	setting.Current()

	var settingEntity entities.WxpaySettingEntity
	copier.Copy(&settingEntity, &setting)

	apiHelpers.ResponseSuccess(c, settingEntity)
}

// @Summary 更新微信支付配置微信支付
// @Produce  json
// @Tags 后台配置管理
// @Param params body entities.WxpaySettingParams true "wxpay_setting params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxpay_setting [put]
// @Security ApiKeyAuth
func UpdateWxpaySetting(c *gin.Context) {
	var err error
	var settingParams entities.WxpaySettingParams
	if err := apiHelpers.ValidateParams(c, &settingParams); err != nil {
		return
	}

	var setting models.WxpaySetting
	setting.Current()
	copier.Copy(&setting, &settingParams)

	err = setting.CreateOrUpdate()
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseSuccess(c, setting)
}

// @Summary 更新微信支付证书
// @Produce  json
// @Tags 后台配置管理
// @Accept  multipart/form-data
// @Param api_client_cert formData file true "wechat pay certification"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/wxpay_setting/cert [post]
// @Security ApiKeyAuth
func UpdateWechatCert(c *gin.Context) {
	var err error
	uploadDir := "./uploads/"
	certName := "apiclient_cert.p12"
	var setting models.WxpaySetting
	setting.Current()

	// header调用Filename方法，就可以得到文件名
	file, _, err := c.Request.FormFile("api_client_cert")

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	err = os.MkdirAll(uploadDir, os.ModePerm)
	filename := filepath.Join(uploadDir, certName)
	out, err := os.Create(filename)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}
	apiHelpers.ResponseOK(c)
}
