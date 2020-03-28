package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/objcoding/wxpay"
	"github.com/xifengzhu/eshop/helpers/setting"
	// "github.com/xifengzhu/eshop/helpers/utils"
	"math"
	"math/rand"
	"time"
)

type OrderStatus int

type Order struct {
	BaseModel

	WxappId            string      `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	OrderNo            string      `gorm:"type: varchar(50); not null; unique_index" json:"order_no"`
	Status             OrderStatus `gorm:"type: tinyint; " json:"status"`
	AddressID          int         `gorm:"-" json:"address_id"`
	ReceiverProperties string      `gorm:"type: varchar(255); " json:"receiver_properties"`
	OuterPayId         string      `gorm:"type: varchar(60); " json:"outer_pay_id"`
	PayAt              *time.Time  `gorm:"type: datetime; " json:"pay_at"`
	ExpressID          int         `gorm:"type: int;" json:"express_id"`
	ExpressAmount      float32     `gorm:"type: decimal(10,2);" json:"express_amount"`
	PayAmount          float32     `gorm:"type: decimal(10,2);" json:"pay_amount"`
	TotalAmount        float32     `gorm:"type: decimal(10,2);" json:"total_amount"`
	UserID             int         `gorm:"type: int; " json:"user_id"`
	BuyerMessage       string      `gorm:"type: varchar(120); " json:"buyer_message"`
	OrderItems         []OrderItem `json:"order_items"`
}

var WxpayClient *wxpay.Client

func (Order) TableName() string {
	return "orders"
}

func init() {
	account := wxpay.NewAccount(setting.WechatAppId, setting.MchID, setting.MchKey, false)
	WxpayClient = wxpay.NewClient(account)

	// 设置http请求超时时间
	WxpayClient.SetHttpConnectTimeoutMs(2000)

	// 设置http读取信息流超时时间
	WxpayClient.SetHttpReadTimeoutMs(1000)

	// 更改签名类型
	WxpayClient.SetSignType("MD5")
}

const (
	wait_buyer_pay OrderStatus = iota
	wait_seller_send_goods
	wait_buyer_confirm_goods
	buyer_confirm_goods
	trade_finished
	canceled
	trade_closed
	refunding
	refunded
)

var OrderStatusEnum = map[string]OrderStatus{
	"wait_buyer_pay":           wait_buyer_pay,
	"wait_seller_send_goods":   wait_seller_send_goods,
	"wait_buyer_confirm_goods": wait_buyer_confirm_goods,
	"buyer_confirm_goods":      buyer_confirm_goods,
	"trade_finished":           trade_finished,
	"canceled":                 canceled,
	"trade_closed":             trade_closed,
	"refunding":                refunding,
	"refunded":                 refunded,
}

var OrderStatusKeys []string = []string{"wait_buyer_pay", "wait_seller_send_goods", "wait_buyer_confirm_goods", "buyer_confirm_goods", "trade_finished", "canceled", "trade_closed", "refunding", "refunded"}

// Scope
func StatusScope(status []int) func(db *gorm.DB) *gorm.DB {
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

//删除数据
func (order Order) Destroy() (err error) {
	err = db.Delete(&order).Error
	return
}

// Callbacks
func (order *Order) BeforeSave() (err error) {
	fmt.Println("=======order before save=============")
	order.setOrderNo()     // 生成订单号
	order.setOrderStatus() // 设置默认订单状态
	return nil
}

func (order *Order) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("=======order after save=============")
	order.caculateExpressAmount(tx)
	order.caculateTotalAmount(tx)

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

func (order *Order) setOrderStatus() {
	order.Status = wait_buyer_pay
}

func (order *Order) generateNo() {
	timeStamp := fmt.Sprintf(time.Now().Format("2006010215040501"))

	randNo := fmt.Sprintf("%02v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100))

	orderNo := timeStamp + randNo
	order.OrderNo = orderNo
}

// TODO: 计算总金额
func (order *Order) caculateTotalAmount(tx *gorm.DB) {
	var totalAmount float32
	for _, orderItem := range order.OrderItems {
		totalAmount += orderItem.TotalAmount
	}
	totalAmount += order.ExpressAmount
	attr := Order{TotalAmount: totalAmount, PayAmount: totalAmount}
	tx.Model(order).Update(attr)
	fmt.Println("=======order after caculateTotalAmount=============", order)
	return
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
func (order *Order) caculateExpressAmount(tx *gorm.DB) {
	var expressAmount float32
	var rules []DeliveryRule
	m := make(map[int][]OrderItem)
	for _, orderItem := range order.OrderItems {
		var goods Goods
		goods.ID = orderItem.GoodsID
		_ = goods.Find()
		rule := goods.DeliveryRule(order.ExpressID, order.AddressID)
		rules = append(rules, rule)
		if _, ok := m[rule.ID]; ok {
			m[rule.ID] = append(m[rule.ID], orderItem)
		}
		m[rule.ID] = []OrderItem{orderItem}
		fmt.Println("===all order_item delivery rules====:", m)
	}
	firstRule := findMaxRuleFirst(rules)
	// 首重/首件 费用
	expressAmount += firstRule.FirstFee
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
			additionalAmount = currentRule.AdditionalFee * float32(math.Ceil(float64(additionWeight)/float64(currentRule.Additional)))
		}
		expressAmount += additionalAmount
		tx.Model(order).Update("express_amount", expressAmount)
		fmt.Println("=======order after caculateExpressAmount=============", order)
		return
	}
}

// Model helper method
func (order Order) OrderNoIsTaken() bool {
	if err := db.Where("order_no = ?", order.OrderNo).First(&order).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

func (order Order) User() (user User) {
	user, _ = GetUserById(order.UserID)
	return user
}

func (order Order) RequestPayment() (map[string]string, error) {
	params := make(wxpay.Params)
	params.SetString("body", "eshop 测试订单").
		SetString("out_trade_no", order.OrderNo).
		SetInt64("total_fee", int64(order.TotalAmount*100)).
		SetString("spbill_create_ip", "127.0.0.1").
		SetString("notify_url", "http://notify.objcoding.com/notify").
		SetString("trade_type", "JSAPI").
		SetString("openid", order.User().OpenId)
	payment, err := WxpayClient.UnifiedOrder(params)
	return payment, err
}

func GetOrderTotal(maps interface{}) (count int) {
	db.Model(&Order{}).Where(maps).Count(&count)
	return
}
