package models

import (
	"errors"
	"fmt"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/objcoding/wxpay"
	"github.com/qor/transition"
	"github.com/xifengzhu/eshop/helpers/export"
	config "github.com/xifengzhu/eshop/initializers"
	"github.com/xifengzhu/eshop/initializers/setting"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Order struct {
	BaseModel

	WxappId            string     `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	OrderNo            string     `gorm:"type: varchar(50); not null; unique_index" json:"order_no"`
	AddressID          int        `gorm:"-" json:"address_id"`
	ReceiverProperties string     `gorm:"type: varchar(250); " json:"receiver_properties"`
	OuterPayId         string     `gorm:"type: varchar(60); " json:"outer_pay_id"`
	PayAt              *time.Time `gorm:"type: datetime; " json:"pay_at"`
	ExpressID          int        `gorm:"type: int;" json:"express_id"`
	ExpressFee         float64    `gorm:"type: decimal(10,2);" json:"express_fee"`
	ProductAmount      float64    `gorm:"type: decimal(10,2);" json:"product_amount"`
	PayAmount          float64    `gorm:"type: decimal(10,2);" json:"pay_amount"`
	AdjustmentAmount   float64    `gorm:"type: decimal(10,2);" json:"adjustment_amount"`
	UserID             int        `gorm:"type: int; " json:"user_id"`
	CouponID           int        `gorm:"type: int; " json:"coupon_id"`
	BuyerMessage       string     `gorm:"type: varchar(120); " json:"buyer_message"`

	OrderItems        []OrderItem   `json:"order_items"`
	User              *User         `gorm:"association_autoupdate:false" json:"user"`
	Coupon            *Coupon       `gorm:"association_autoupdate:false" json:"coupon"`
	Express           *Express      `gorm:"association_autoupdate:false" json:"express,omitempty"`
	Address           *Address      `gorm:"association_autoupdate:false" json:"address,omitempty"`
	LatestPaymentTime *time.Time    `gorm:"type: datetime;" json:"latest_payment_time"`
	AllAdjustments    []*Adjustment `json:"all_adjustments,omitempty"`
	DirectAdjustments []*Adjustment `gorm:"direct_polymorphic:Target;" json:"adjustments,omitempty"`
	transition.Transition
}

var (
	WxpayClient *wxpay.Client
	OrderFSM    = transition.New(&Order{})
	RemainTime  = time.Minute * 2
	PaidStates  = []string{"wait_seller_send_goods", "wait_buyer_confirm_goods", "buyer_confirm_goods", "trade_finished", "refunding", "refunded"}
)

func init() {
	initWechatAccount()
	defineState()
}

func initWechatAccount() {
	account := wxpay.NewAccount(setting.WechatAppId, setting.MchID, setting.MchKey, false)

	// 设置证书
	account.SetCertData("./uploads/apiclient_cert.p12")

	// new client
	WxpayClient = wxpay.NewClient(account)

	// 设置http请求超时时间
	WxpayClient.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	WxpayClient.SetHttpReadTimeoutMs(1000)

	// 更改签名类型
	WxpayClient.SetSignType("MD5")

	WxpayClient.SetAccount(account)
}

func defineState() {
	// Define Order's States
	OrderFSM.Initial("wait_buyer_pay")
	OrderFSM.State("wait_seller_send_goods")
	OrderFSM.State("wait_buyer_confirm_goods")
	OrderFSM.State("buyer_confirm_goods")
	OrderFSM.State("trade_finished")
	OrderFSM.State("canceled")
	OrderFSM.State("trade_closed")
	OrderFSM.State("refunding")
	OrderFSM.State("refunded")

	// Define State Event
	OrderFSM.Event("paid").
		To("wait_seller_send_goods").
		From("wait_buyer_pay").
		After(func(order interface{}, tx *gorm.DB) error {
			// TOTO: 1. 根据减库存规则减库存
			log.Println("after paid....")
			return nil
		})
	OrderFSM.Event("ship").
		To("wait_buyer_confirm_goods").
		From("wait_seller_send_goods").
		After(func(order interface{}, tx *gorm.DB) error {
			// TOTO: 1. 发货提醒 2. 7天之后自动确认收货任务
			log.Println("after ship....")
			return nil
		})
	OrderFSM.Event("confirm").
		To("buyer_confirm_goods").
		From("wait_buyer_confirm_goods").
		After(func(order interface{}, tx *gorm.DB) error {
			// TOTO: 1. 确认收货之后，7天自动交易完成，关闭售后
			log.Println("after confirm....")
			return nil
		})
	OrderFSM.Event("finish").To("trade_finished").From("buyer_confirm_goods")
	OrderFSM.Event("cancel").To("canceled").From("wait_buyer_pay")
	OrderFSM.Event("refund").To("refunding").From("wait_seller_send_goods")
	OrderFSM.Event("drawback").To("refunded").From("refunding")
	OrderFSM.Event("close").To("trade_closed").From("wait_buyer_pay")

}

// Callbacks
func (order *Order) BeforeCreate() (err error) {
	fmt.Println("=======order before create=============")
	order.setOrderNo()    // 生成订单号
	order.setOrderState() // 设置默认订单状态
	order.setExpressFee()
	order.setProductAmount()
	order.applyCoupon()
	order.setPayAmount()
	return
}

func (order *Order) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("=======order AfterCreate =============", order.PayAmount)
	order.setLatestPaymentTime(tx)
	order.lockCoupon(tx)
	order.enqueueCloseOrderJob()
	// order.removeFromCartItem()
	err = order.reductStock() // if reduceType == 0
	return
}

func (order *Order) reductStock() (err error) {
	for _, orderItem := range order.OrderItems {
		err = orderItem.ReductStock()
		if err != nil {
			break
		}
	}
	return
}

func (order *Order) RestoreStock() {
	// restore stock number
	for _, orderItem := range order.OrderItems {
		orderItem.RestoreStock()
	}
}

func (order *Order) removeFromCartItem() {
	db.Where("user_id =? AND checked IS TRUE", order.UserID).Delete(&CarItem{})
}

func (order *Order) PreOrder() {
	order.setExpressFee()
	order.setProductAmount()
	order.applyOrSetDefaultCoupon()
	order.setPayAmount()
}

func (order Order) prepareCouponOrder() (corder CouponOrder) {
	corder = CouponOrder{
		ProductAmount: order.ProductAmount,
		FreightFee:    order.ExpressFee,
	}
	for _, item := range order.OrderItems {
		a := CouponOrderItem{
			ProductAmount: item.ProductAmount,
			ResourceID:    item.GoodsID,
			ResourceType:  "goods",
		}
		corder.Items = append(corder.Items, a)
	}
	return
}

func (order *Order) applyOrSetDefaultCoupon() {
	if order.CouponID == 0 {
		coupon := order.biggestDiscountCoupons()
		order.CouponID = coupon.ID
		order.Coupon = coupon
	}
	order.applyCoupon()
}

func (order *Order) applyCoupon() {
	coupon, err := order.currentCoupon()
	if err == nil {
		corder := order.prepareCouponOrder()
		result := coupon.Apply(corder)
		for _, oi := range result.Items {
			for i, item := range order.OrderItems {
				if oi.ResourceID == item.GoodsID {
					adjustment := &Adjustment{
						Amount:     oi.ReduceAmount * -1,
						Label:      fmt.Sprintf("Coupon: %s", coupon.Name),
						SourceType: "Coupon",
						SourceID:   order.CouponID,
					}
					order.AllAdjustments = append(order.AllAdjustments, adjustment)
					order.OrderItems[i].Adjustments = append(order.OrderItems[i].Adjustments, adjustment)
					order.OrderItems[i].AdjustmentAmount += adjustment.Amount
					order.OrderItems[i].TotalAmount = order.OrderItems[i].ProductAmount + adjustment.Amount
					order.AdjustmentAmount += adjustment.Amount
					break
				}
			}
		}
	}
}

// Private method
func (order *Order) setOrderNo() {
	for {
		order.generateNo()
		if !order.OrderNoIsTaken() {
			break
		}
	}
}

func (order *Order) setLatestPaymentTime(tx *gorm.DB) {
	DeadLine := order.CreatedAt.Add(RemainTime)
	tx.Model(order).Updates(Order{LatestPaymentTime: &DeadLine})
	return
}

func (order *Order) currentCoupon() (coupon Coupon, err error) {
	if order.CouponID != 0 {
		coupon.ID = order.CouponID
		err = Find(&coupon, Options{})
	} else {
		err = errors.New("not have coupon")
	}
	return
}

func (order *Order) lockCoupon(tx *gorm.DB) (err error) {
	log.Println("========lock coupon id ======", order.CouponID)
	coupon, err := order.currentCoupon()
	if err == nil {
		if coupon.State != "actived" {
			err = errors.New("优惠券已失效！")
		} else {
			tx.Model(coupon).Updates(Coupon{State: "lock"})
		}
	}
	return
}

func (order *Order) setOrderState() {
	order.State = "wait_buyer_pay"
}

func (order *Order) generateNo() {
	timeStamp := fmt.Sprintf(time.Now().Format("2006010215040501"))

	randNo := fmt.Sprintf("%02v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100))

	orderNo := timeStamp + randNo
	order.OrderNo = orderNo
}

func (order *Order) setPayAmount() {
	payAmount := order.ExpressFee + order.ProductAmount + order.AdjustmentAmount
	order.PayAmount = payAmount
}

// 计算商品金额
func (order *Order) setProductAmount() {
	var productAmount float64
	for i, _ := range order.OrderItems {
		order.OrderItems[i].CaculateProductAmount()
		productAmount += order.OrderItems[i].ProductAmount
	}
	order.ProductAmount = productAmount
}

// 可用的优惠券列表
func (order Order) AvaliableCoupons() []Coupon {
	avaliableCoupons := []Coupon{}
	corder := order.prepareCouponOrder()
	coupons := order.User.GetActivedCoupons()
	for _, coupon := range coupons {
		result := coupon.Apply(corder)
		if result.ReduceAmount > 0 || result.ReduceFreight > 0 {
			avaliableCoupons = append(avaliableCoupons, coupon)
		}
	}
	return avaliableCoupons
}

// 找优惠力度最大的优惠券
func (order *Order) biggestDiscountCoupons() *Coupon {
	var maxReduceAmount float64
	var maxDiscountCoupon Coupon
	corder := order.prepareCouponOrder()
	coupons := order.User.GetActivedCoupons()
	for _, coupon := range coupons {
		result := coupon.Apply(corder)
		reduceTotal := result.ReduceAmount + result.ReduceFreight
		if reduceTotal > maxReduceAmount {
			maxDiscountCoupon = coupon
			maxReduceAmount = reduceTotal
		}
	}
	return &maxDiscountCoupon
}

// 前提条件
// 商品 A 按件计费 首件 1  运费  8   续重 1  运费 3   运费模板A
// 商品 B 按件计费 首件 1  运费  8   续重 1  运费 3   运费模板A
// 商品 c 按件计费 首件 1  运费  20  续重 1  运费 10  运费模板B
// 商品 D 按重量计 首重 1  运费  12  续重 1  运费 5   运费模板C

// case 1
// 当我同时购买2件商品A和2件商品B，因为这两款商品都使用同一个运费模板，则只使用其中一款商品的首件运费，其余商品直接按照续件的运费进行计算，即购买2件商品A和2件商品B时，系统计算的运费应该为： 8元+3元+3元+3元=17元

// case 2
// 当我同时购买2件商品A和2件商品C，这两款商品都是按照件数计算运费，但使用的是不同的运费模板，则淘宝会比较这两款商品中首件的运费，选择首件运费最大的费用作为首件费用（商品C首件运费20元），然后忽略商品A的首件运费，商品A全部按照该商品的续件费用进行计算运费。系统计算的运费应该为： 20元+10元+3元+3元=36元

// case 3
// 当我同时购买2件商品A和2kg商品D，商品A按照件数计算运费，商品D按照重量计算运费，则淘宝会比较这两款商品中首件（首公斤）的运费，选择首件（首公斤）运费最大的费用作为首件（首公斤）费用（商品D首公斤运费12元），然后忽略商品A的首件运费，商品A全部按照该商品的续件费用进行计算运费。系统计算的运费应该为： 12元+5元+3元+3元=23元
func (order *Order) setExpressFee() {
	var expressFee float64
	eligible := order.eligibleShopFreeFreightSetting()
	if !eligible {
		expressFee = order.freightWithDeliveryRule()
	}
	order.ExpressFee = expressFee
}

// TODO:
func (order *Order) eligibleShopFreeFreightSetting() bool {
	var setting GlobalSetting
	setting.Current()
	return order.ProductAmount >= setting.FreeFreightAmount
}

func (order *Order) freightWithDeliveryRule() float64 {
	var expressFee float64
	var rules []DeliveryRule
	m := make(map[int][]OrderItem)
	for _, orderItem := range order.OrderItems {
		var goods Goods
		goods.ID = orderItem.GoodsID
		Find(&goods, Options{})
		rule := goods.DeliveryRule(order.ExpressID, order.AddressID)
		rules = append(rules, rule)
		if _, ok := m[rule.ID]; ok {
			m[rule.ID] = append(m[rule.ID], orderItem)
		}
		m[rule.ID] = []OrderItem{orderItem}
	}
	firstRule := findMaxRuleFirst(rules)
	// 首重/首件 费用
	expressFee += firstRule.FirstFee
	var additionalAmount float64
	for ruleID, orderItems := range m {
		var currentRule DeliveryRule
		for _, rule := range rules {
			if rule.ID == ruleID {
				currentRule = rule
				break
			}
		}
		if currentRule.Delivery.IsCaculateByNum() {
			var totalNum int
			for _, orderItem := range orderItems {
				totalNum += orderItem.TotalNum
			}
			var additionNum int
			if currentRule.ID == firstRule.ID {
				additionNum = totalNum - int(firstRule.First)
			} else {
				additionNum = totalNum
			}
			additionalAmount = currentRule.AdditionalFee * float64(math.Ceil(float64(additionNum)/float64(currentRule.Additional)))
		} else {
			var totalWeight float64
			for _, orderItem := range orderItems {
				totalWeight += orderItem.GoodsWeight
			}
			var additionWeight float64
			if currentRule.ID == firstRule.ID {
				additionWeight = totalWeight - firstRule.First
			} else {
				additionWeight = totalWeight
			}
			if currentRule.Additional != 0 {
				additionalAmount = currentRule.AdditionalFee * float64(math.Ceil(float64(additionWeight)/float64(currentRule.Additional)))
			} else {
				additionalAmount = 0
			}
		}
	}
	expressFee += additionalAmount
	return expressFee
}

func findMaxRuleFirst(rules []DeliveryRule) (max DeliveryRule) {
	max = rules[0]
	for _, rule := range rules {
		if rule.First > max.First {
			max = rule
		}
	}
	return max
}

// Model helper method
func (order Order) OrderNoIsTaken() bool {
	if err := db.Where("order_no = ?", order.OrderNo).First(&order).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func (order Order) enqueueCloseOrderJob() {
	config.JobEnqueuer.EnqueueIn("close_order", 60*15, work.Q{"order_id": order.ID})
}

func (order Order) RequestPayment(ip string) (map[string]string, error) {
	payment, err := order.UnifiedOrder(ip)
	timastamp := strconv.FormatInt(time.Now().Unix(), 10)
	signParams := make(wxpay.Params)
	signParams.SetString("package", "prepay_id="+payment["prepay_id"]).
		SetString("nonceStr", payment["nonce_str"]).
		SetString("timeStamp", timastamp).
		SetString("appId", setting.WechatAppId).
		SetString("signType", "MD5")

	paymentParams := map[string]string{
		"timeStamp": timastamp,
		"nonceStr":  payment["nonce_str"],
		"package":   "prepay_id=" + payment["prepay_id"],
		"signType":  "MD5",
		"paySign":   WxpayClient.Sign(signParams),
	}

	log.Println("=====paymentParams=======", paymentParams)

	return paymentParams, err
}

func (order Order) UnifiedOrder(ip string) (map[string]string, error) {
	params := make(wxpay.Params)
	params.SetString("body", "eshop 测试订单").
		SetString("out_trade_no", order.OrderNo).
		SetInt64("total_fee", int64(order.PayAmount*100)).
		SetString("spbill_create_ip", ip).
		SetString("notify_url", "http://notify.objcoding.com/notify").
		SetString("trade_type", "JSAPI").
		SetString("openid", order.User.OpenId)

	payment, err := WxpayClient.UnifiedOrder(params)

	log.Println("=====payment=======", payment)
	if err != nil {
		log.Println(err)
	}
	return payment, err
}

func (order Order) Close() (err error) {
	tx := db.Begin()
	if err = OrderFSM.Trigger("close", &order, db, "auto close order after 15min"); err != nil {
		log.Println("====close order failed===", err)
		return err
	}

	if err = tx.Set("gorm:association_autoupdate", false).Model(&order).UpdateColumn("state", order.State).Error; err != nil {
		tx.Rollback()
		return err
	}
	order.RestoreStock()
	// Or commit the transaction
	tx.Commit()
	return

}

func (order Order) Pay() (err error) {
	tx := db.Begin()
	if err = OrderFSM.Trigger("paid", &order, tx, "admin user pay order"); err != nil {
		log.Println("====pay order failed===", err)
		return err
	}

	if err = tx.Set("gorm:association_autoupdate", false).Model(&order).UpdateColumn("state", order.State).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(order.Coupon).Updates(Coupon{State: "used"})
	// Or commit the transaction
	tx.Commit()
	return
}

func (order Order) Ship(tx *gorm.DB) (err error) {
	if err = OrderFSM.Trigger("ship", &order, tx, "admin user ship order"); err != nil {
		log.Println("====ship order failed===", err)
		return err
	}
	if err = tx.Set("gorm:association_autoupdate", false).Model(&order).UpdateColumn("state", order.State).Error; err != nil {
		return err
	}
	return nil
}

func (order Order) Finish() (err error) {
	tx := db.Begin()
	if err = OrderFSM.Trigger("finish", &order, tx, "system finish order"); err != nil {
		log.Println("====ship order failed===", err)
		return err
	}
	if err = tx.Set("gorm:association_autoupdate", false).Model(&order).UpdateColumn("state", order.State).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (order Order) Confirm(operator string) (err error) {
	tx := db.Begin()
	message := fmt.Sprintf("%s finish order", operator)
	if err = OrderFSM.Trigger("confirm", &order, tx, message); err != nil {
		log.Println("====ship order failed===", err)
		return err
	}
	if err = tx.Set("gorm:association_autoupdate", false).Model(&order).UpdateColumn("state", order.State).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (order Order) DestroyOrderItems() {
	db.Where("order_id = ?", order.ID).Delete(OrderItem{})
}

func (order Order) OrderItemsCount() (total int) {
	type Result struct{ Count int }
	var result Result
	db.Table("order_item").Where("order_id = ?", order.ID).Select("sum(total_num) as count").Scan(&result)
	return result.Count
}

func (order Order) Export() (string, error) {
	var orders []Order
	All(&orders, Options{})
	data := [][]string{}
	titles := []string{"编号", "订单号", "件数", "产品金额", "运费", "支付金额", "收货信息", "支付时间"}
	for _, v := range orders {
		var payAt string
		if v.PayAt == nil {
			payAt = ""
		} else {
			payAt = v.PayAt.Format("2006-01-02 03:04:05 PM")
		}
		values := []string{
			strconv.Itoa(v.ID),
			v.OrderNo,
			strconv.Itoa(v.OrderItemsCount()),
			fmt.Sprintf("%f", v.ProductAmount),
			fmt.Sprintf("%f", v.ExpressFee),
			fmt.Sprintf("%f", v.PayAmount),
			v.ReceiverProperties,
			payAt,
		}
		data = append(data, values)
	}

	return export.Exec("订单列表", titles, data)
}
