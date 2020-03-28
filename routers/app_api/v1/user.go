package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"strconv"

	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/helpers/wechat"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// Binding from JSON
type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// @Summary 通过当前登录用户信息
// @Produce  json
// @Tags 用户
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/users/mine [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	apiHelpers.ResponseSuccess(c, currentUser)
}

// @Summary 编辑用户
// @Produce  json
// @Tags 用户
// @Param params body UserInfo true "user info"
// @Success 200 {object} apiHelpers.Response
// @Router /app_api/v1/users/mine [put]
// @Security ApiKeyAuth
func EditUser(c *gin.Context) {
	var userInfo UserInfo
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)
	if err := c.BindJSON(&userInfo); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	} else {
		data := utils.StructToMap(userInfo)
		models.EditUser(int(user.ID), data)
		newUser, _ := models.GetUserById(int(user.ID))
		apiHelpers.ResponseSuccess(c, newUser)
	}
}

// @Summary 用户通过微信的code获取auth token
// @Tags 用户
// @Description 用户获取token
// @Accept  json
// @Produce  json
// @Param code query string true "wechat code"
// @Success 200 {object} apiHelpers.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /app_api/v1/user/auth [get]
func AuthWithWechat(c *gin.Context) {
	code := c.Query("code")
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
		apiHelpers.ResponseSuccess(c, data)
	} else {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err)
	}
}

// @Summary 通过用户IDauth token
// @Tags 用户
// @Description 用户获取token
// @Accept  json
// @Produce  json
// @Param user_id query integer true "user id"
// @Success 200 {object} apiHelpers.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /app_api/v1/user/fake_token [get]
func GetToken(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	params := make(map[string]interface{})
	params["id"] = userID
	params["resource"] = "user"
	token := utils.Encode(params)
	apiHelpers.ResponseSuccess(c, token)
}

// @Summary 用户解码token/仅限测试使用
// @Produce  json
// @Tags 用户
// @Param token query string true "Token"
// @Success 200 {object} apiHelpers.Response
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
		apiHelpers.ResponseSuccess(c, result)
	}
}
