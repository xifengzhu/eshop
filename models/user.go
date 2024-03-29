package models

import (
	"github.com/jinzhu/gorm"
	// "github.com/xifengzhu/eshop/helpers/utils"
	"errors"
)

type User struct {
	BaseModel

	Username  string    `json:"username"`
	Gender    int       `json:"gender"`
	Avatar    string    `json:"avatar"`
	OpenId    string    `gorm:"type: varchar(50); unique_index; not null"json:"open_id"`
	Addresses []Address `json:"addresses"`
	CarItems  []CarItem `json:"shopping_cart_items"`
	Orders    []Order   `json:"orders"`
}

func FindOrCreateUserByOpenId(openId string) (user User) {
	db.FirstOrCreate(&user, User{OpenId: openId})
	return
}

func GetUserById(id int) (user User, err error) {
	err = db.First(&user, id).Error
	return user, nil
}

func GetUserByOpenId(openId string) (interface{}, error) {
	var user User
	if err := db.Where("open_id = ?", openId).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return user, nil
}

func GetUserTotal(maps interface{}) (count int) {
	db.Model(&User{}).Where(maps).Count(&count)
	return
}

func ExistUserByOpenId(open_id string) bool {
	var user User
	db.Select("id").Where("open_id = ?", open_id).First(&user)
	if user.ID > 0 {
		return true
	}
	return false
}

func ExistUserByID(id int) bool {
	var user User
	db.Select("id").Where("id = ?", id).First(&user)
	if user.ID > 0 {
		return true
	}

	return false
}

func AddUser(username string, avatar string, open_id string) User {
	user := &User{
		Username: username,
		Avatar:   avatar,
		OpenId:   open_id,
	}
	db.Create(user)
	return *user
}

// =======  address ============
func (user User) GetAddresses() (addresses []Address) {
	db.Set("gorm:auto_preload", true).Model(&user).Association("Addresses").Find(&addresses)
	return
}

func (user User) GetAddressByID(addressID int) (address Address, err error) {
	err = db.Set("gorm:auto_preload", true).Find(&address, "user_id = ? AND id = ?", user.ID, addressID).Error
	return
}

// =======  shopping cart ============
func (user User) GetShoppingCartItems() (cartItems []CarItem) {
	db.Set("gorm:auto_preload", true).Model(&user).Association("CarItems").Find(&cartItems)
	return
}

func (user User) GetShoppingCartItemByID(itemID int) (cartItem CarItem, err error) {
	err = db.Set("gorm:auto_preload", true).Find(&cartItem, "id = ? AND user_id = ?", itemID, user.ID).Error
	return
}

func (user User) GetShoppingCartItemByIDs(itemIDs []int) (cartItems []CarItem, err error) {
	err = db.Set("gorm:auto_preload", true).Find(&cartItems, "user_id = ? AND id IN (?)", user.ID, itemIDs).Error
	return
}

func (user User) GetCheckedShoppingCartItems() (cartItems []CarItem, err error) {
	err = db.Set("gorm:auto_preload", true).Find(&cartItems, "user_id = ? AND checked IS TRUE", user.ID).Error
	return
}

func (user User) FindShoppingCartItemByGoodsID(goodsID int) (cartItem CarItem, err error) {
	err = db.Where("user_id = ? AND goods_id = ?", user.ID, goodsID).First(&cartItem).Error
	return
}

func (user User) CatchCoupon(couponTemplateID int) (bool, error) {
	var template CouponTemplate
	template.ID = couponTemplateID
	err := Find(&template, Options{})
	if err != nil {
		return false, errors.New("优惠券不存在~")
	}
	if template.Stock <= template.CouponsCount {
		return false, errors.New("优惠券已经领完~")
	}
	if template.CatchLimit > template.CatchCountForUser(user.ID) {
		coupon := template.GenerateCouponsData(1)[0]
		coupon.UserID = user.ID
		db.Create(&coupon)
		return true, nil
	} else {
		return false, errors.New("您已超过限制领取次数！")
	}
}

func (user User) GetActivedCoupons() (coupons []Coupon) {
	db.Model(&Coupon{}).Where("user_id = ? AND state = 'actived'", user.ID).Find(&coupons)
	return
}

// =======  order ============
func (user User) GetOrder(orderID int) (order Order, err error) {
	err = db.Preload("OrderItems").Preload("User").Where("id = ? AND user_id = ? ", orderID, user.ID).First(&order).Error
	return
}
