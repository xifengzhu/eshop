package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/initializers/wechat"
	"github.com/xifengzhu/eshop/models"
	. "github.com/xifengzhu/eshop/routers/app_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 通过当前登录用户信息
// @Produce  json
// @Tags 用户
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/users/mine [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	currentUser, _ := c.Get("resource")
	ResponseSuccess(c, currentUser)
}

// @Summary 编辑用户
// @Produce  json
// @Tags 用户
// @Param params body params.UserInfo true "user info"
// @Success 200 {object} helpers.Response
// @Router /app_api/v1/users/mine [put]
// @Security ApiKeyAuth
func EditUser(c *gin.Context) {
	var userInfo UserInfo
	currentUser, _ := c.Get("resource")
	user := currentUser.(models.User)
	if err := ValidateParams(c, &userInfo, "json"); err != nil {
		return
	}

	changedAttrs := models.User{}
	copier.Copy(&changedAttrs, &userInfo)

	models.Update(&user, &changedAttrs)

	ResponseSuccess(c, user)
}

// @Summary 用户通过微信的code获取auth token
// @Tags 用户
// @Description 用户获取token
// @Accept  json
// @Produce  json
// @Param params body params.AuthParams true "wechat auth params"
// @Param code query string true "wechat code"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /app_api/v1/user/auth [post]
func AuthWithWechat(c *gin.Context) {
	var auth AuthParams
	if err := ValidateParams(c, &auth, "json"); err != nil {
		return
	}

	code := auth.Code
	result, err := wechat.CodeToSession(code)
	openId := result["openid"].(string)
	var data interface{}
	if err == nil {
		if openId != "" {
			user := models.FindOrCreateUserByOpenId(openId)
			params := make(map[string]interface{})
			params["id"] = user.ID
			params["resource"] = "user"
			token := utils.Encode(params)
			data = token
		} else {
			data = result
		}
		ResponseSuccess(c, data)
	} else {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
	}
}

// @Summary 通过用户IDauth token
// @Tags 用户
// @Description 用户获取token
// @Accept  json
// @Produce  json
// @Param resource_id query integer true "resource id"
// @Param resource_type query string true "resource type"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /app_api/v1/user/fake_token [get]
func GetToken(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("resource_id"))
	params := make(map[string]interface{})
	params["id"] = userID
	params["resource"] = c.Query("resource_type")
	token := utils.Encode(params)
	ResponseSuccess(c, token)
}

// @Summary 用户解码token/仅限测试使用
// @Produce  json
// @Tags 用户
// @Param token query string true "Token"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /app_api/v1/user/verify [post]
// @Security ApiKeyAuth
func VerifyToken(c *gin.Context) {
	token := c.Query("token")
	result, err := utils.Decode(token)
	if err != nil {
		fmt.Println("err", err)
	} else {
		ResponseSuccess(c, result)
	}
}
