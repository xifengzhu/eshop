package models

import (
	"strconv"
	"strings"

	"github.com/xifengzhu/eshop/helpers/utils"
)

type DeliveryRule struct {
	BaseModel

	WxappId       string    `gorm:"type: varchar(50); not null" json:"wxapp_id"`
	First         float32   `gorm:"type: decimal(10,2); " json:"first"` // 首件/首重
	FirstFee      float32   `gorm:"type: decimal(10,2); " json:"first_fee"`
	Additional    float32   `gorm:"type: decimal(10,2);" json:"additional"`     // 续件/续重
	AdditionalFee float32   `gorm:"type: decimal(10,2);" json:"additional_fee"` // 续件/续重
	Region        string    `sql:"type: json;" json:"region,omitempty"`         // 可配送区域(省id集)
	DeliveryID    int       `gorm:"type: int;" json:"delivery_id"`
	ExpressID     int       `gorm:"type: int; not null" json:"express_id"`
	Delivery      *Delivery `json:"delivery,omitempty"`
	Destroy       bool      `sql:"-" json:"_destroy,omitempty"`
}

func (rule DeliveryRule) SuitableProvinceIDs() (provinceIDs []int) {
	ids := strings.Split(rule.Region, ",")
	for _, id := range ids {
		intID, _ := strconv.Atoi(id)
		provinceIDs = append(provinceIDs, intID)
	}
	return
}

func (rule DeliveryRule) SuitableProvince() (provinces []string) {
	ids := rule.SuitableProvinceIDs()
	db.Select("name").Where("id = ?", ids).Find(&provinces)
	return
}

func (rule DeliveryRule) Hit(provinceId int) bool {
	ids := rule.SuitableProvinceIDs()
	return utils.ContainsInt(ids, provinceId)
}
