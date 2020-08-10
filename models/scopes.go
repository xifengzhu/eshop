package models

import (
	"github.com/jinzhu/gorm"
	"github.com/xifengzhu/eshop/helpers/utils"
	"time"
)

type Amount struct {
	Total float64
}

type DateGroupCount struct {
	Date string `json:"date"`
	Data int    `json:"data"`
}

type DateGroupAmount struct {
	Date string  `json:"date"`
	Data float64 `json:"data"`
}

// ===== common scopes =======
func StateScope(state []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("state in (?)", state)
	}
}

func RecentDay(days int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 计算begining
		beginOfToday := utils.BeginningOfDay(time.Now())
		h, _ := time.ParseDuration("-24h")
		h1 := beginOfToday.Add(time.Duration(days-1) * h)
		return db.Where("created_at > ?", h1)
	}
}

// ===== order scopes =======
func TodayNewOrderCount() (count int) {
	db.Model(&Order{}).Scopes(RecentDay(1), StateScope(PaidStates)).Count(&count)
	return
}

func Recent7DayNewOrderCount() (count int) {
	db.Model(&Order{}).Scopes(RecentDay(7), StateScope([]string{"wait_seller_send_goods"})).Count(&count)
	return
}

func TotalTodayOrderAmount() float64 {
	var amount Amount
	db.Model(&Order{}).Scopes(RecentDay(1), StateScope(PaidStates)).Select("sum(pay_amount) as total").Scan(&amount)
	return amount.Total
}

func Total7DayOrderAmount() float64 {
	var amount Amount
	db.Model(&Order{}).Scopes(RecentDay(7), StateScope(PaidStates)).Select("sum(pay_amount) as total").Scan(&amount)
	return amount.Total
}

func TotalOrderAmount() float64 {
	var amount Amount
	db.Model(&Order{}).Scopes(StateScope(PaidStates)).Select("sum(pay_amount) as total").Scan(&amount)
	return amount.Total
}

func OrderNumTrend() []DateGroupCount {
	data := make([]DateGroupCount, 30)
	db.Model(&Order{}).Scopes(RecentDay(30)).Select("date(created_at) as date, count(*) as data").Group("date(created_at)").Scan(&data)
	return data
}

func PaidOrderNumTrend() []DateGroupCount {
	data := make([]DateGroupCount, 30)
	db.Model(&Order{}).Scopes(RecentDay(30), StateScope(PaidStates)).Select("date(created_at) as date, count(*) as data").Group("date(created_at)").Scan(&data)
	return data
}

func ActualOrderAmountTrend() []DateGroupAmount {
	data := make([]DateGroupAmount, 30)
	db.Model(&Order{}).Scopes(RecentDay(30)).Select("date(created_at) as date, sum(product_amount) as data").Group("date(created_at)").Scan(&data)
	return data
}

func PaidOrderAmountTrend() []DateGroupAmount {
	data := make([]DateGroupAmount, 30)
	db.Model(&Order{}).Scopes(RecentDay(30), StateScope(PaidStates)).Select("date(created_at) as date, sum(pay_amount) as data").Group("date(created_at)").Scan(&data)
	return data
}

// ===== user scopes =======
func TotalUserCount() (count int) {
	db.Model(&User{}).Count(&count)
	return
}

func TodayNewUserCount() (count int) {
	db.Model(&User{}).Scopes(RecentDay(1)).Count(&count)
	return
}
