package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/helpers"
)

type PolicyRuleParams struct {
	RoleName string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

// @Summary 添加权限
// @Produce  json
// @Tags 后台角色权限
// @Param params body PolicyRuleParams true "权限"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/cabin_rules [post]
// @Security ApiKeyAuth
func AddPolicy(c *gin.Context) {
	var rule PolicyRuleParams
	if err := apiHelpers.ValidateParams(c, &rule, "json"); err != nil {
		return
	}

	res, _ := models.Enforcer.AddPolicy(rule.RoleName, rule.Path, rule.Method)
	if !res {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "添加失败，权限已经存在")
		return
	}
	apiHelpers.ResponseOK(c)
}

// @Summary 删除权限
// @Produce  json
// @Tags 后台角色权限
// @Param params body PolicyRuleParams true "权限"
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/cabin_rules [delete]
// @Security ApiKeyAuth
func RemovePolicy(c *gin.Context) {
	var rule PolicyRuleParams
	if err := c.ShouldBindJSON(&rule); err != nil {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, err.Error())
	}
	res, _ := models.Enforcer.RemovePolicy(rule.RoleName, rule.Path, rule.Method)
	if !res {
		apiHelpers.ResponseError(c, e.INVALID_PARAMS, "删除失败，权限不存在")
		return
	}
	apiHelpers.ResponseOK(c)
}
