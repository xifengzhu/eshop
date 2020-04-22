package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/qor/transition"
	"github.com/xifengzhu/eshop/helpers/setting"
)

var db *gorm.DB

type BaseModel struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT" form:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"default: CURRENT_TIMESTAMP" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default: CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" form:"updated_at"`
}

func init() {
	var err error

	tablePrefix := setting.DatabaseTablePrefix

	db, err = gorm.Open(setting.DatabaseType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseUser,
		setting.DatabasePassword,
		setting.DatabaseHost,
		setting.DatabasePort,
		setting.DatabaseName))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	migration()
}

func migration() {
	//Migrate the schema
	db.AutoMigrate(&User{}, &Tag{}, &Address{}, &AppSetting{}, &CarItem{}, &Category{}, &City{}, &Delivery{}, &DeliveryRule{}, &Goods{}, &Logistic{}, &Order{}, &OrderItem{}, &Product{}, &PropertyName{}, &PropertyValue{}, &Province{}, &Region{}, &Tag{}, &User{}, &WebPage{}, &WxappPage{}, &Express{}, &AdminUser{}, &transition.StateChangeLog{}, &ProductGroup{})
}

func CloseDB() {
	defer db.Close()
}
