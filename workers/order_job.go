package workers

import (
	"github.com/gocraft/work"
	"github.com/xifengzhu/eshop/models"

	"log"
)

func (c *Context) CloseOrder(job *work.Job) error {
	orderID := job.ArgInt64("order_id")
	log.Println("========CloseOrder=====", orderID)
	var order models.Order
	order.ID = int(orderID)
	models.Find(&order, models.Options{Preloads: []string{"OrderItems"}})
	if order.State == "wait_buyer_pay" {
		order.Close()
	}

	return nil
}

func (c *Context) ConfirmOrder(job *work.Job) error {
	orderID := job.ArgInt64("order_id")
	var order models.Order
	order.ID = int(orderID)
	models.Find(&order, models.Options{})
	if order.State == "wait_buyer_confirm_goods" {
		order.Confirm("system")
	}

	return nil
}
