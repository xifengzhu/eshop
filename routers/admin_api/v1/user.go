package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	. "github.com/xifengzhu/eshop/models"
	_ "github.com/xifengzhu/eshop/routers/admin_api/params"
	. "github.com/xifengzhu/eshop/routers/helpers"
	"strconv"
)

// @Summary 获取用户列表
// @Produce  json
// @Tags 后台用户管理
// @Param params query params.QueryUserParams true "query params"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/users [get]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	pagination := SetDefaultPagination(c)

	var model User
	result := &[]User{}

	Search(&model, &SearchParams{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := Collection{Pagination: pagination, List: result}

	ResponseSuccess(c, response)
}

// @Summary 获取用户详情
// @Produce  json
// @Tags 后台用户管理
// @Param id path int true "user id"
// @Success 200 {object} helpers.Response
// @Router /admin_api/v1/users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	var user User
	user.ID, _ = strconv.Atoi(c.Param("id"))

	err := Find(&user, Options{Preloads: []string{"Orders"}})
	if err != nil {
		ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	ResponseSuccess(c, user)
}
