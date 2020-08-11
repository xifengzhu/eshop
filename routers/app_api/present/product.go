package present

import (
	"github.com/xifengzhu/eshop/models"
	"time"
)

type ProductDetailEntity struct {
	ID              int                    `json:"id"`
	Cover           string                 `json:"cover"`
	MainPictures    models.JSON            `json:"main_pictures"`
	Name            string                 `json:"name"`
	Content         string                 `json:"content"`
	DeductStockType int                    `json:"deduct_stock_type"`
	SalesInitial    int                    `json:"sales_initial"`
	SalesActual     int                    `json:"sales_actual"`
	Position        int                    `json:"position"`
	Price           float64                `json:"price"`
	IsOnline        bool                   `json:"is_online"`
	DeletedAt       *time.Time             `json:"deleted_at"`
	DeliveryID      int                    `json:"delivery_id"`
	CategoryID      int                    `json:"category_id"`
	Goodses         []models.Goods         `json:"goodses"`
	Specifications  []models.Specification `json:"specifications"`
}

type ProductEntity struct {
	ID    int    `json:"id"`
	Cover string `json:"cover"`
	// MainPictures models.JSON `json:"main_pictures"`
	Name         string     `json:"name"`
	Content      string     `json:"content"`
	SalesInitial int        `json:"sales_initial"`
	SalesActual  int        `json:"sales_actual"`
	Position     int        `json:"position"`
	Price        float64    `json:"price"`
	IsOnline     bool       `json:"is_online"`
	DeletedAt    *time.Time `json:"deleted_at"`
	CategoryID   int        `json:"category_id"`
}
