package params

type QrCodeParams struct {
	Page      string `form:"page"`
	Scene     string `form:"scene"`
	Width     int    `form:"width" `
	IsHyaline bool   `form:"is_hyaline" `
	Binary    bool   `form:"binary" `
}
