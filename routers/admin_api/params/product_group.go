package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
	"time"
)

type ProductGroupParams struct {
	Name       string      `json:"name"`
	Remark     string      `json:"remark"`
	ProductIDs models.JSON `json:"product_ids"`
	Key        string      `json:"key"`
}

type QueryProductGroupParams struct {
	utils.Pagination
	Name            string    `json:"q[name]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
}
