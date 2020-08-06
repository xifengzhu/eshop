package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/models"
)

func CurrentAdmin(c *gin.Context) (admin models.AdminUser) {
	admin, _ := c.MustGet("resource").(models.AdminUser)
	return
}
