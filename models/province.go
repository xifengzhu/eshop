package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Province struct {
	BaseModel

	Name   string  `gorm:"type: varchar(64); not null" json:"name"`
	Cities []*City `json:"cities"`
}

type City struct {
	BaseModel

	Name       string    `gorm:"type: varchar(64); not null" json:"name"`
	ProvinceID int       `gorm:"type: int; not null" json:"province_id"`
	Province   *Province `json:"province,omitempty"`
	Regions    []*Region `json:"regions"`
}

type Region struct {
	BaseModel

	Name   string `gorm:"type: varchar(64); not null" json:"name"`
	CityID int    `gorm:"type: int; not null" json:"city_id"`
	City   *City  `json:"city,omitempty"`
}

func init() {
	var province Province
	if province.isRecordExist() {
		log.Println("=====Province: seed data is exist====")
	} else {
		file, err := ioutil.ReadFile("./seeds/data/address.json")

		if err != nil {
			log.Println("Province: load address seed json error: ", err)
		}

		var provinces []Province
		err = json.Unmarshal(file, &provinces)

		if err != nil {
			log.Println("Province: json unmarshal error: ", err)
		}

		for _, pro := range provinces {
			db.Create(&pro)
		}
	}
}

func (Province) isRecordExist() bool {
	var count int
	db.Model(&Province{}).Count(&count)
	return count > 0
}
