package v1

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/initializers/setting"
	"github.com/xifengzhu/eshop/initializers/wechat"
	. "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"io/ioutil"
	"net/http"
	"os"
)

// @Summary 获取微信二维码
// @Produce  json
// @Tags 微信接口
// @Param params query params.QrCodeParams true "二维码参数"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/wechat/wxacode [get]
func GetWxaCode(c *gin.Context) {

	var params QrCodeParams
	if err := ValidateParams(c, &params, "query"); err != nil {
		return
	}

	// wechat.RefreshWxAccessToken()

	response, err := wechat.GetWxaCodeUnLimit(params.Page, params.Scene, params.Width, params.IsHyaline)
	wechatResp, err := ioutil.ReadAll(response.Body)

	var errResp map[string]interface{}
	json.Unmarshal(wechatResp, &errResp)

	defer response.Body.Close()

	_, ok := errResp["errcode"]
	if err != nil || ok {
		ResponseError(c, e.WECHAT_QRCCODE_ERROR, "生成二维码失败！")
		return
	}

	if params.Binary {
		c.Data(http.StatusOK, "application/octet-stream", wechatResp)
		return
	} else {
		// 写入文件
		has := md5.Sum(wechatResp)
		fileName := fmt.Sprintf("%x.jpg", has)
		file := fmt.Sprintf("./public/qrcode/%s", fileName)
		f, _ := os.Create(file)
		defer f.Close()
		f.Write(wechatResp)
		resp := map[string]string{"path": getQrcodeFullUrl(fileName)}
		ResponseSuccess(c, resp)
		return
	}
	ResponseError(c, e.WECHAT_QRCCODE_ERROR, "生成二维码失败！")
}

func getQrcodeFullUrl(name string) string {
	return setting.Domain + "/" + setting.QrcodeSavePath + name
}
