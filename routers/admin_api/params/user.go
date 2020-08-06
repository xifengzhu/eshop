package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type QueryUserParams struct {
	utils.Pagination
	OpenId          string     `json:"q[open_id_eq]"`
	Username        string     `json:"q[username_cont]"`
	Created_at_gteq *time.Time `json:"q[created_at_gteq]"`
	Created_at_lteq *time.Time `json:"q[created_at_lteq]"`
}
