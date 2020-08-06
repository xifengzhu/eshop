package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type QueryCouponTemplateParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}

type Config struct {
	MinAmount    float64    `json:"min_amount" validate:"required,gte=0"`
	ResourceType string     `json:"resource_type,omitempty"`
	Resources    []int      `json:"resources,omitempty" validate:"required_with=ResourceType"`
	DateType     string     `json:"date_type" validate:"required,oneof='fix_term' 'time_range'"`
	FixTerm      int        `json:"fix_term,omitempty" validate:"rfe=DateType:fix_term"`
	StartAt      *time.Time `json:"start_at,omitempty" validate:"required_with=EndAt,rfe=DateType:time_range"`
	EndAt        *time.Time `json:"end_at,omitempty" validate:"required_with=StartAt,rfe=DateType:time_range""`

	ReduceAmount float64 `json:"reduce_amount,omitempty" validate:"required_without=Percentage"`
	Percentage   int     `json:"percentage,omitempty" validate:"required_without=ReduceAmount"`
}

type CouponTemplateParams struct {
	Code       string     `json:"code"`
	Name       string     `json:"name" validate:"required"`
	Kind       string     `json:"kind" validate:"required,oneof='fixed_amount' 'percentage'"`
	Creator    string     `json:"creator"`
	Stock      int        `json:"stock" validate:"required"`
	CatchLimit int        `json:"catch_limit" validate:"required"`
	StartAt    *time.Time `json:"start_at" validate:"required"`
	EndAt      *time.Time `json:"end_at" validate:"required,gtfield=StartAt"`
	Configs    Config     `json:"configs" validate:"required,dive"`
}
