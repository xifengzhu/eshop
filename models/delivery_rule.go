package models

import (
	"encoding/json"
	"github.com/xifengzhu/eshop/helpers/utils"
	"strings"
)

type DeliveryRule struct {
	BaseModel

	First         float64   `gorm:"type: decimal(10,2); " json:"first"` // 首件/首重
	FirstFee      float64   `gorm:"type: decimal(10,2); " json:"first_fee"`
	Additional    float64   `gorm:"type: decimal(10,2);" json:"additional"`     // 续件/续重
	AdditionalFee float64   `gorm:"type: decimal(10,2);" json:"additional_fee"` // 续件/续重
	Region        JSON      `sql:"type: json;" json:"region,omitempty"`         // 可配送区域(省id集)
	DeliveryID    int       `gorm:"type: int;" json:"delivery_id"`
	Delivery      *Delivery `json:"delivery,omitempty"`
	Destroy       bool      `sql:"-" json:"_destroy,omitempty"`
	Position      int       `gorm:"type: int; " json:"position"`
	RegionName    string    `sql:"-" json:"region_names,omitempty"`
}

func (rule *DeliveryRule) AfterFind() (err error) {
	rule.SetRegionNames()
	return
}

func (rule *DeliveryRule) SetRegionNames() {
	var names []string
	regions := rule.regionIDs()
	if len(regions) > 0 {
		db.Model(&Province{}).Where("id IN (?)", regions).Pluck("name", &names)
		rule.RegionName = strings.Join(names[:], ",")
	}
}

func (rule DeliveryRule) Hit(provinceId int) bool {
	regions := rule.regionIDs()
	return utils.ContainsInt(regions, provinceId)
}

func (rule DeliveryRule) regionIDs() (regions []int) {
	if rule.Region == nil {
		return
	}
	if err := json.Unmarshal(rule.Region, &regions); err != nil {
		panic(err)
	}
	return
}
