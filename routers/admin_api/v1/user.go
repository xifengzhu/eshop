package v1

import (
	"github.com/gin-gonic/gin"
)

type QueryUserParams struct {
}

// @Summary 获取用户列表
// @Produce  json
// @Tags 后台用户管理
// @Param params query QueryUserParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/users [get]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {

}

// @Summary 获取用户详情
// @Produce  json
// @Tags 后台用户管理
// @Param id path int true "user id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {

}
