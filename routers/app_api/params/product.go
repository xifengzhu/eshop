package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type ProductQueryParams struct {
	utils.Pagination
	CategoryName string `form:"category_name"`
	Keyword      string `form:"keyword" validate:"required_without=CategoryName"`
}
