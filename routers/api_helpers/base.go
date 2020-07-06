package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/unknwon/com"
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

func ValidateParams(c *gin.Context, params interface{}) error {

	if err := c.ShouldBindJSON(params); err != nil {
		ResponseError(c, e.INVALID_PARAMS, err)
		return err
	}
	validate := validator.New()
	errs := validate.Struct(params)
	if errs != nil {
		ResponseError(c, e.INVALID_PARAMS, errs)
		return errs
	}

	return nil
}

func SetDefaultPagination(c *gin.Context) (pagination *utils.Pagination) {
	perPage := com.StrTo(c.DefaultQuery("per_page", "10")).MustInt()
	page := com.StrTo(c.DefaultQuery("page", "1")).MustInt()
	Sort := c.DefaultQuery("order_by", "id desc")
	pagination = &utils.Pagination{Page: page, PerPage: perPage, Sort: Sort}
	return
}

func ResponseError(c *gin.Context, code int, err error) {
	response := &Response{Code: code, Msg: err.Error(), Data: nil}
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
