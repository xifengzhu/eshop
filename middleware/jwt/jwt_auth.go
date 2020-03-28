package jwt

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("[middleware] request Authorization: ", c.GetHeader("Authorization"))
		token := c.GetHeader("Authorization")
		if token == "" {
			Unauthorized(c)
		} else {
			claims, err := utils.Decode(token)
			if err != nil {
				Unauthorized(c)
				return
			}
			// find current user
			if claims["id"] == nil || claims["resource"] == nil {
				Unauthorized(c)
			} else {
				resourceID := int(claims["id"].(float64))
				legalResourceTyle := []string{"user", "admin"}
				legal := utils.ContainsString(legalResourceTyle, claims["resource"].(string))
				if legal {
					var err error
					var resource interface{}
					if claims["resource"] == "user" {
						resource, err = models.GetUserById(resourceID)
						log.Print("[middleware] currentUser==: ", resource)
					}
					if claims["resource"] == "admin" {
						resource, err = models.GetAdminUserById(resourceID)
						log.Print("[middleware] currentAdmin==: ", resource)
					}
					if err != nil {
						Unauthorized(c)
						return
					}
					// 继续交由下一个路由处理,并将解析出的信息传递下去
					c.Set("claims", claims)
					c.Set("resource", resource)
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"code": 400,
						"msg":  "token错误",
					})
				}
			}
		}
	}
}

func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code": 401,
		"msg":  "无效的token，无权限访问",
	})
}
