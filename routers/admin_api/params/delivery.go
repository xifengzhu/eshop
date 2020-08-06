package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
)

type DeliveryRuleParams struct {
	ID            int         `json:"id,omitempty"`
	First         float64     `json:"first" ` // 首件/首重
	FirstFee      float64     `json:"first_fee" `
	Additional    float64     `json:"additional" `     // 续件/续重
	AdditionalFee float64     `json:"additional_fee" ` // 续件/续重
	Region        models.JSON `json:"region" `         // 可配送区域(省id集)
	Destroy       bool        `json:"_destroy,omitempty"`
	Position      int         `json:"position"`
}

type DeliveryParams struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty" validate:"required"`
	Way           int                  `json:"way" validate:"required"` // 1 为按件计费 2 按重量计费
	DeliveryRules []DeliveryRuleParams `json:"delivery_rules" validate:"required,dive"`
}

type QueryDeliveryParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}
