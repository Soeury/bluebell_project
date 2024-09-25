package controllers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 定义一个全局的翻译器 , 将参数校验失败的错误信息翻译一下返回给前端
var trans ut.Translator

// InitTrans初始化翻译器
func InitTrans(locale string) (err error) {

	// 修改 gin 框架中的validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个方法，让校验失败返回的错误信息中的字段名字改为自定义的而不是结构体中定义的
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		// 第一个参数值备用的语言环境
		// 后面的参数是应该支持的语言环境(可以是多个)
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'accept-language'
		var ok bool
		// 也可以使用 uni.FindTranslater(...) 传入多个 locale 查找
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// 移除校验失败返回的错误信息中的结构体前缀
// "xxx.oo" : "www"
// 移除后:
// "oo" : "www"
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
