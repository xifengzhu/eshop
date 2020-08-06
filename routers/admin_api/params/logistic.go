package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type QueryLogisticParams struct {
	utils.Pagination
	ExpressCompany string `json:"q[express_company_cont]"`
	ExpressNo      string `json:"q[express_no_cont]"`
}
