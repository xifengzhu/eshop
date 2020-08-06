package params

import (
	"github.com/xifengzhu/eshop/helpers/utils"
	"github.com/xifengzhu/eshop/models"
)

type GoodsParams struct {
	ID         int     `json:"id,omitempty"`
	Name       string  `json:"name"`
	Properties string  `json:"properties"`
	Image      string  `json:"image"`
	SkuNo      string  `json:"sku_no"`
	StockNum   int     `json:"stock_num"`
	Position   int     `json:"position"`
	Price      float64 `json:"price"`
	LinePrice  float64 `json:"line_price"`
	Weight     float64 `json:"weight"`
	Destroy    bool    `json:"_destroy,omitempty"`
}

type ProductParams struct {
	Name            string        `json:"name"`
	Content         string        `json:"content"`
	DeductStockType int           `json:"deduct_stock_type"`
	SalesInitial    int           `json:"sales_initial"`
	Position        int           `json:"position"`
	Price           float64       `json:"price"`
	IsOnline        bool          `json:"is_online"`
	DeliveryID      int           `json:"delivery_id"`
	Goodses         []GoodsParams `json:"goodses"`
	CategoryIDs     []int         `json:"category_ids"`
	SerialNo        string        `json:"serial_no"`
	MainPictures    models.JSON   `json:"main_pictures"`
	ShareDesc       string        `json:"share_desc"`
	ShareCover      string        `json:"share_cover"`
}

type QueryProductParams struct {
	utils.Pagination
	IsOnline   bool    `json:"q[is_online_eq]"`
	Name       string  `json:"q[name_cont]"`
	Price_gteq float64 `json:"q[price_gteq]"`
	Price_lteq float64 `json:"q[price_lteq]"`
}
