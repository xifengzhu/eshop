package params

type WebPageParams struct {
	Title   string `json:"title" validate:"required"`
	Cover   string `json:"cover"`
	Content string `json:"content" validate:"required"`
}
