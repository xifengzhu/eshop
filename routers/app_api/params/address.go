package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
)

type AddressParams struct {
	UserID    int    `json:"user_id"`
	Region    string `json:"region" validate:"required"`
	Province  string `json:"province" validate:"required"`
	City      string `json:"city" validate:"required"`
	Detail    string `json:"detail" validate:"required"`
	isDefault bool   `json:"is_default" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Receiver  string `json:"receiver" validate:"required"`
}

type AddressQueryParams struct {
	utils.Pagination
}
