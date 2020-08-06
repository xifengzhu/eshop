package apiHelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xifengzhu/eshop/helpers/e"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/xifengzhu/eshop/routers/validators"
)

var (
	translator *ut.UniversalTranslator
	validate   *validator.Validate
)

func init() {
	zh := zh.New()
	translator = ut.New(zh)
	validate = validator.New()
	validators.RegisterCustomValidations()
}

func ValidateParams(c *gin.Context, params interface{}, bindType string) (err error) {

	locale := c.DefaultQuery("locale", "zh")
	translator, _ := translator.GetTranslator(locale)
	if locale == "zh" {
		zh_translations.RegisterDefaultTranslations(validate, translator)
	} else {
		en_translations.RegisterDefaultTranslations(validate, translator)
	}

	if bindType == "json" {
		err = c.ShouldBindJSON(params)
	} else {
		err = c.ShouldBindQuery(params)
	}

	if err != nil {
		ResponseError(c, e.INVALID_PARAMS, err.Error())
		return
	}

	err = validate.Struct(params)
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
