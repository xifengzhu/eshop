package entities

import (
	"time"
)

type WxpaySettingParams struct {
	WxappId       string `json:"wxapp_id,omitempty" `
	AppName       string `json:"app_name,omitempty"`
	AppSecret     string `json:"app_secret,omitempty"`
	Mchid         string `json:"mchid,omitempty"`
	Apikey        string `json:"apikey,omitempty"`
	NotifyUrl     string `json:"notify_url,omitempty"`
	ApiClientCert string `json:"api_client_cert,omitempty"`
}

type WxpaySettingEntity struct {
	WxappId   string     `json:"wxapp_id"`
	AppName   string     ` json:"app_name"`
	Mchid     string     ` json:"mchid"`
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
