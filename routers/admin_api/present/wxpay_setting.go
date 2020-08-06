package present

import (
	"time"
)

type WxpaySettingEntity struct {
	WxappId   string     `json:"wxapp_id"`
	AppName   string     ` json:"app_name"`
	Mchid     string     ` json:"mchid"`
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
