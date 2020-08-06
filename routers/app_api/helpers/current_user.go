package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/models"
)

func CurrentUser(c *gin.Context) models.User {
	return c.MustGet("resource").(models.User)
}
