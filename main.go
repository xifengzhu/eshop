package main

import (
	"fmt"
	"github.com/xifengzhu/eshop/helpers/setting"
	"github.com/xifengzhu/eshop/routers"
	_ "github.com/xifengzhu/eshop/workers"
	"net/http"
)

func main() {
	// @title Eshop API
	// @version 1.0
	// @description This is a sample server eshop server.

	// @host 127.0.0.1:8000
	// @BasePath /

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
