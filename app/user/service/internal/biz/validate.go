package biz

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	v1 "github.com/user-center/user-center-backend/api/user/service/v1"
	"reflect"
)

type ValidateUseCase struct {
	trans    ut.Translator
	validate *validator.Validate
}

func NewValidateUseCase() *ValidateUseCase {
	zh_ch := zh.New()
	uni := ut.New(zh_ch)                // 万能翻译器，保存所有的语言环境和翻译数据
	trans, _ := uni.GetTranslator("zh") // 翻译器
	Validate := validator.New()
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
	_ = zh_translations.RegisterDefaultTranslations(Validate, trans)
	return &ValidateUseCase{
		validate: Validate,
		trans:    trans,
	}
}

func (u *ValidateUseCase) ParamsValidate(object interface{}) error {
	err := u.validate.Struct(object)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			return v1.ErrorValidateError(e.Translate(u.trans))
		}
	}
	return nil
}
