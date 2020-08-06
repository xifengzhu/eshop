package params

type PropertyValueParams struct {
	ID             int    `json:"id,omitempty"`
	Value          string `json:"value" validate:"required"`
	PropertyNameID int    `json:"property_name_id" validate:"required"`
	Destroy        bool   `json:"_destroy,omitempty"`
}
