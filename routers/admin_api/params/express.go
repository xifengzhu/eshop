package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type ExpressParams struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type QueryExpressParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
	Code string `json:"q[code_cont]"`
}
