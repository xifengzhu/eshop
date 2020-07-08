package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
	"strconv"
	"time"
)

type QueryUserParams struct {
	utils.Pagination
	OpenId          string     `json:"q[open_id_eq]"`
	Username        string     `json:"q[username_cont]"`
	Created_at_gteq *time.Time `json:"q[created_at_gteq]"`
	Created_at_lteq *time.Time `json:"q[created_at_lteq]"`
}

// @Summary 获取用户列表
// @Produce  json
// @Tags 后台用户管理
// @Param params query QueryUserParams true "query params"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/users [get]
// @Security ApiKeyAuth
func GetUsers(c *gin.Context) {
	pagination := apiHelpers.SetDefaultPagination(c)

	var model models.User
	result := &[]models.User{}

	models.Search(&model, &Search{Pagination: pagination, Conditions: c.QueryMap("q")}, &result)

	response := apiHelpers.Collection{Pagination: pagination, List: result}

	apiHelpers.ResponseSuccess(c, response)
}

// @Summary 获取用户详情
// @Produce  json
// @Tags 后台用户管理
// @Param id path int true "user id"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	var user models.User
	user.ID, _ = strconv.Atoi(c.Param("id"))

	err := models.Find(&user, Query{Preloads: []string{"Orders"}})
	if err != nil {
		apiHelpers.ResponseError(c, e.ERROR_NOT_EXIST, err.Error())
		return
	}

	apiHelpers.ResponseSuccess(c, user)
}
