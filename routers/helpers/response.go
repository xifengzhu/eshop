package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg"`
}

type Collection struct {
	Pagination *utils.Pagination `json:"pagination,omitempty"`
	List       interface{}       `json:"list"`
}

func ResponseError(c *gin.Context, code int, errMsg string) {
	response := &Response{Code: code, Msg: errMsg}
	c.AbortWithStatusJSON(http.StatusBadRequest, response)
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	response := &Response{Code: e.SUCCESS, Data: data, Msg: e.GetMsg(e.SUCCESS)}
	c.JSON(http.StatusOK, response)
}

func ResponseOK(c *gin.Context) {
	response := &Response{Code: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)}
	c.JSON(http.StatusOK, response)
}
