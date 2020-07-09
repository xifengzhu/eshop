package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/unknwon/com"
	"github.com/xifengzhu/eshop/helpers/e"
	"github.com/xifengzhu/eshop/helpers/utils"
	"net/http"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Translator *ut.UniversalTranslator
	Validate   *validator.Validate
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

func init() {
	zh := zh.New()
	Translator = ut.New(zh)
	Validate = validator.New()
}

func ValidateParams(c *gin.Context, params interface{}) (err error) {

	locale := c.DefaultQuery("locale", "zh")
	translator, _ := Translator.GetTranslator(locale)
	if locale == "zh" {
		zh_translations.RegisterDefaultTranslations(Validate, translator)
	} else {
		en_translations.RegisterDefaultTranslations(Validate, translator)
	}

	err = c.BindJSON(params)
	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	err = Validate.Struct(params)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		var errMsg string
		for _, e := range errs {
			errMsg = e.Translate(translator)
			break
		}
		ResponseError(c, e.INVALID_PARAMS, errMsg)
	}
	return
}

func SetDefaultPagination(c *gin.Context) (pagination *utils.Pagination) {
	perPage := com.StrTo(c.DefaultQuery("per_page", "10")).MustInt()
	page := com.StrTo(c.DefaultQuery("page", "1")).MustInt()
	Sort := c.DefaultQuery("order_by", "id desc")
	pagination = &utils.Pagination{Page: page, PerPage: perPage, Sort: Sort}
	return
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
