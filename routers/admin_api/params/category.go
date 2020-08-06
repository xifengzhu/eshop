package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type CategoryParams struct {
	Name     string `json:"name" validate:"required" `
	Position int    `json:"position"`
	ParentID int    `json:"parent_id"`
	Image    string `json:"image"`
}

type QueryCategoryParams struct {
	utils.Pagination
	Name            string    `json:"q[name]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	ParentID        int       `json:"q[parent_id_eq]"`
}
