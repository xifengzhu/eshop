package ip_filter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/utils"
)

var BlockList = []string{}

func IPFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if utils.ContainsString(BlockList, ip) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "Permission denied",
			})
		}
	}
}
