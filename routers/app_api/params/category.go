package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type CategoryProductQueryParams struct {
	utils.Pagination
	CategoryID string `json:"cagegory_id"`
}
