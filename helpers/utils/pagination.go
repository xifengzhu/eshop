package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xifengzhu/eshop/initializers/setting"
)

type Pagination struct {
	Page    int    `form:"page" json:"page" validate:"gte=0"`
	PerPage int    `form:"per_page" json:"per_page" validate:"lt=100"`
	Sort    string `form:"sort" json:"sort"`
	Total   int    `form:"total" json:"total"`
}

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.PageSize
	}

	return result
}
