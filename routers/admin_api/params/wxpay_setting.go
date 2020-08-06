package params

type WxpaySettingParams struct {
	WxappId       string `json:"wxapp_id,omitempty" `
	AppName       string `json:"app_name,omitempty"`
	AppSecret     string `json:"app_secret,omitempty"`
	Mchid         string `json:"mchid,omitempty"`
	Apikey        string `json:"apikey,omitempty"`
	NotifyUrl     string `json:"notify_url,omitempty"`
	ApiClientCert string `json:"api_client_cert,omitempty"`
}
