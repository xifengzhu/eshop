package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type PropertyParams struct {
	ID             int            `json:"id,omitempty"`
	Name           string         `json:"name" validate:"required"`
	PropertyValues []ValuesParams `json:"property_values"  validate:"dive"`
}

type ValuesParams struct {
	ID             int    `json:"id,omitempty"`
	Value          string `json:"value" validate:"required"`
	PropertyNameID int    `json:"property_name_id"`
	Destroy        bool   `json:"_destroy,omitempty"`
}

type QueryPropertyParams struct {
	utils.Pagination
	Name string `json:"q[name_cont]"`
}
