package params

type CartItemParams struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
	GoodsID  int `json:"goods_id" validate:"required"`
}

type CartItemIDParams struct {
	ItemIDs []int `json:"item_ids" validate:"required,gt=0"`
}

type CartItemQtyParams struct {
	ItemID   int `json:"item_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}
