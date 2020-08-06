package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type AddAdminUserParams struct {
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
	Password string `json:"password"  validate:"required,gte=6,lt=12"`
}

type UpdateAdminUserParams struct {
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Status   string `json:"status,omitempty" validate:"oneof=active banned"`
}

type AddRolesParams struct {
	RoleIds []int `json:"role_ids"`
}

type QueryAdminUserParams struct {
	utils.Pagination
	Email_cont      string    `json:"q[email_cont]"`
	Created_at_gteq time.Time `json:"q[created_at_gteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Created_at_lteq time.Time `json:"q[created_at_lteq]" time_format:"2006-01-02T15:04:05Z07:00"`
	Status          []int     `json:"q[status_in]"`
}
