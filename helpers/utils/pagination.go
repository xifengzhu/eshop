package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"github.com/xifengzhu/eshop/helpers/setting"
)

type Pagination struct {
	Page    int    `json:"page" binding:"gte=0"`
	PerPage int    `json:"per_page" binding:"lt=100"`
	Sort    string `json:"sort"`
	Total   int    `json:"total"`
}

func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.PageSize
	}

	return result
}
