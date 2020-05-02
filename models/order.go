package models

import (
	"fmt"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/objcoding/wxpay"
	"github.com/qor/transition"
	"github.com/xifengzhu/eshop/helpers/setting"
	"github.com/xifengzhu/eshop/helpers/utils"
	"log"
	"math"
	"math/rand"
	"time"
)

type Order struct {
	BaseModel

	WxappId            string      `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	OrderNo            string      `gorm:"type: varchar(50); not null; unique_index" json:"order_no"`
	AddressID          int         `gorm:"-" json:"address_id"`
	ReceiverProperties string      `gorm:"type: varchar(255); " json:"receiver_properties"`
	OuterPayId         string      `gorm:"type: varchar(60); " json:"outer_pay_id"`
	PayAt              *time.Time  `gorm:"type: datetime; " json:"pay_at"`
	ExpressID          int         `gorm:"type: int;" json:"express_id"`
	ExpressFee         float32     `gorm:"type: decimal(10,2);" json:"express_fee"`
	ProductAmount      float32     `gorm:"type: decimal(10,2);" json:"product_amount"`
	PayAmount          float32     `gorm:"type: decimal(10,2);" json:"pay_amount"`
	UserID             int         `gorm:"type: int; " json:"user_id"`
	BuyerMessage       string      `gorm:"type: varchar(120); " json:"buyer_message"`
	OrderItems         []OrderItem `json:"order_items"`
	User               *User       `json:"user"`
	Express            *Express    `json:"express,omitempty"`
	transition.Transition
}

func (Order) TableName() string {
	return "orders"
}

var (
	WxpayClient *wxpay.Client
	OrderFSM    = transition.New(&Order{})
)

func init() {
	initWechatAccount()
	defineState()
}

func initWechatAccount() {
	account := wxpay.NewAccount(setting.WechatAppId, setting.MchID, setting.MchKey, false)
	WxpayClient = wxpay.NewClient(account)

	// 设置证书
	// WxpayClient.certData = AppSetting.ApiClientCertData

	// 设置http请求超时时间
	WxpayClient.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	WxpayClient.SetHttpReadTimeoutMs(1000)

	// 更改签名类型
	WxpayClient.SetSignType("MD5")
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
	OrderFSM.Event("cancel").To("wait_buyer_pay").From("canceled")
	OrderFSM.Event("refund").To("refunding").From("wait_seller_send_goods")
	OrderFSM.Event("drawback").To("refunded").From("refunding")
	OrderFSM.Event("close").To("wait_buyer_pay").From("trade_closed")

}

// Scope
func StateScope(status []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status in (?)", status)
	}
}

// DB
func (order *Order) Create() (err error) {
	err = db.Set("gorm:save_associations", true).Create(&order).Error
	if err != nil {
		return
	}
	db.Set("gorm:auto_preload", true).Find(&order)
	return
}

// Callbacks
func (order *Order) BeforeCreate() (err error) {
	fmt.Println("=======order before save=============")
	order.setOrderNo()    // 生成订单号
	order.setOrderState() // 设置默认订单状态
	return nil
}

func (order *Order) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("=======order after save=============")
	order.caculateExpressFee(tx)
	order.caculateProductAmount(tx)
	order.caculatePayAmount(tx)
	order.enqueueCloseOrderJob()

	err = reductStock()
	if err != nil {
		return
	}

	err = removeFromCartItem()
	if err != nil {
		return
	}
	return
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

// TODO: 减库存
func reductStock() (err error) {
	return
}

// TODO: 从购物车移除
func removeFromCartItem() (err error) {
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

// 计算支付金额
func (order *Order) caculatePayAmount(tx *gorm.DB) {
	var totalAmount float32
	totalAmount = order.ExpressFee + order.ProductAmount
	tx.Model(order).Updates(Order{PayAmount: totalAmount})
	return
}

// 计算商品金额
func (order *Order) caculateProductAmount(tx *gorm.DB) {
	var productAmount float32
	for _, orderItem := range order.OrderItems {
		productAmount += orderItem.TotalAmount
	}
	tx.Model(order).Updates(Order{ProductAmount: productAmount})
	return
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
func (order *Order) caculateExpressFee(tx *gorm.DB) {
	var expressFee float32
	eligible := order.eligibleShopFreeFreightSetting()
	if eligible {
		expressFee = 0
	} else {
		expressFee = order.freightWithDeliveryRule()
	}
	tx.Model(order).Update("express_fee", expressFee)
	return
}

// TODO:
func (order *Order) eligibleShopFreeFreightSetting() bool {
	var setting GlobalSetting
	setting.Current()
	return order.ProductAmount >= setting.FreeFreightAmount
}

func (order *Order) freightWithDeliveryRule() float32 {
	var expressFee float32
	var rules []DeliveryRule
	m := make(map[int][]OrderItem)
	for _, orderItem := range order.OrderItems {
		var goods Goods
		goods.ID = orderItem.GoodsID
		FindResource(&goods, Options{})
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
	var additionalAmount float32
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
			additionalAmount = currentRule.AdditionalFee * float32(math.Ceil(float64(additionNum)/float64(currentRule.Additional)))
		} else {
			var totalWeight float32
			for _, orderItem := range orderItems {
				totalWeight += orderItem.GoodsWeight
			}
			var additionWeight float32
			if currentRule.ID == firstRule.ID {
				additionWeight = totalWeight - firstRule.First
			} else {
				additionWeight = totalWeight
			}
			if currentRule.Additional != 0 {
				additionalAmount = currentRule.AdditionalFee * float32(math.Ceil(float64(additionWeight)/float64(currentRule.Additional)))
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
	enqueuer := work.NewEnqueuer(setting.RedisNamespace, utils.RedisPool)
	enqueuer.EnqueueIn("close_order", 900, work.Q{"order_id": order.ID})
}

// func (order Order) User() (user User) {
// 	user, _ = GetUserById(order.UserID)
// 	return user
// }

func (order Order) RequestPayment() (map[string]string, error) {
	params := make(wxpay.Params)
	params.SetString("body", "eshop 测试订单").
		SetString("out_trade_no", order.OrderNo).
		SetInt64("total_fee", int64(order.ProductAmount*100)).
		SetString("spbill_create_ip", "127.0.0.1").
		SetString("notify_url", "http://notify.objcoding.com/notify").
		SetString("trade_type", "JSAPI").
		SetString("openid", order.User.OpenId)
	payment, err := WxpayClient.UnifiedOrder(params)
	return payment, err
}

func (order Order) Close() (err error) {
	tx := db.Begin()
	if err = OrderFSM.Trigger("close", &order, db, "auto close order after 15min"); err != nil {
		log.Println("====close order failed===", err)
		return err
	}

	if err = tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
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

	if err = tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Or commit the transaction
	tx.Commit()
	return
}

func (order Order) Ship(tx *gorm.DB) (err error) {
	if err = OrderFSM.Trigger("ship", &order, tx, "admin user pay order"); err != nil {
		log.Println("====ship order failed===", err)
		return err
	}
	if err = tx.Save(&order).Error; err != nil {
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
	if err = tx.Save(&order).Error; err != nil {
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
	if err = db.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (order Order) DestroyOrderItems() {
	db.Where("order_id = ?", order.ID).Delete(OrderItem{})
}
