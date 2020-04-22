package v1

import (
	// "errors"
	// // "fmt"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/xifengzhu/eshop/helpers/e"
	// "github.com/xifengzhu/eshop/helpers/utils"
	// "github.com/xifengzhu/eshop/models"
	apiHelpers "github.com/xifengzhu/eshop/routers/api_helpers"
)

// @Summary dashboard
// @Produce  json
// @Tags 数据统计
// @Success 200 {object} apiHelpers.Response
// @Router /admin_api/v1/dashboard [get]
func Dashboard(c *gin.Context) {
	apiHelpers.ResponseSuccess(c, nil)
}
