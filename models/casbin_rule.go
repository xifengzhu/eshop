package models

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"log"
)

var Enforcer *casbin.Enforcer

func init() {
	fmt.Println("init cabins rule")
	adapter, _ := gormadapter.NewAdapterByDB(db)
	Enforcer, _ = casbin.NewEnforcer("conf/rbac_model.conf", adapter)
	Enforcer.EnableLog(true)
	err := Enforcer.LoadPolicy()
	if err != nil {
		log.Println("loadPolicy error")
		panic(err)
	}
}
