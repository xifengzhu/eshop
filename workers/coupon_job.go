package workers

import (
	"github.com/gocraft/work"
	"github.com/xifengzhu/eshop/models"
	"log"
)

func (c *Context) ExpireCoupon(job *work.Job) error {
	CouponID := job.ArgInt64("order_id")
	log.Println("========expire coupon=====", CouponID)
	var coupon models.Coupon
	coupon.ID = int(CouponID)
	models.Find(&coupon, models.Options{})
	if coupon.State != "expired" {
		coupon.Expire()
	}

	return nil
}
