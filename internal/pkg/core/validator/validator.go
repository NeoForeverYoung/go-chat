package validator

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	_ = Initialize()
}

func Initialize() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		chinese := zh.New()
		uni := ut.New(chinese)
		trans, _ = uni.GetTranslator("zh")

		// 注册一个函数，获取struct tag里自定义的label作为字段名
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("label")
			return name
		})

		registerCustomValidator(v, trans)

		return zhTranslations.RegisterDefaultTranslations(v, trans)
	}

	return nil
}

// TODO 这里不是很理解
// Translate 翻译错误信息
func Translate(err error) string {
	// 检查错误是否是验证器错误
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		// 遍历验证器错误，获取每个错误的信息
		for _, err := range errs {
			// 使用翻译器翻译错误信息
			return err.Translate(trans)
		}
	}

	// 如果错误不是验证器错误，直接返回错误信息
	return err.Error()
}

func Validate(value interface{}) error {
	return binding.Validator.Engine().(*validator.Validate).Struct(value)
}

// registerCustomValidator 注册自定义验证器
func registerCustomValidator(v *validator.Validate, trans ut.Translator) {
	phone(v, trans)
	ids(v, trans)
}
