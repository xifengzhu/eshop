package models

import (
	"time"
)

type Menu struct {
	BaseModel

	Name       string     `json:"name"`
	Position   int        `json:"position"`
	State      int        `json:"state"` // 1:启用 2:禁用
	ShowStatus int        `json:"state"` // 1:显示 2:隐藏
	Remark     string     `json:"remark"`
	Icon       string     `json:"icon"`
	SubMenu    []*Menu    `gorm:"foreignkey:parent_id" json:"submenu,omitempty"`
	Parent     *Menu      `json:"parent,omitempty"`
	Creator    string     `json:"creator"`
	DeletedAt  *time.Time `gorm:"type: datetime; " json:"deleted_at"`
	Path       string     `json:path`
	// Roles      []*Role    `gorm:"many2many:role_menu;"`
}
