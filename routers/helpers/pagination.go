package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xifengzhu/eshop/helpers/utils"
)

func SetDefaultPagination(c *gin.Context) (pagination *utils.Pagination) {
	perPage := com.StrTo(c.DefaultQuery("per_page", "10")).MustInt()
	page := com.StrTo(c.DefaultQuery("page", "1")).MustInt()
	Sort := c.DefaultQuery("order_by", "id asc")
	pagination = &utils.Pagination{Page: page, PerPage: perPage, Sort: Sort}
	return
}
