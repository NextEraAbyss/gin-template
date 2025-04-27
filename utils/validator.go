package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans ut.Translator
)

// 初始化验证器
func init() {
	// 获取gin的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册结构体验证tag
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// 注册中文翻译器
		zhTrans := zh.New()
		uni := ut.New(zhTrans, zhTrans)
		trans, _ = uni.GetTranslator("zh")
		_ = zhtranslations.RegisterDefaultTranslations(v, trans)

		// 注册自定义验证器
		registerCustomValidators(v)
	}
}

// 注册自定义验证器
func registerCustomValidators(v *validator.Validate) {
	// 示例：注册一个手机号验证器
	_ = v.RegisterValidation("mobile", validateMobile)
}

// 手机号验证函数
func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 简单的中国手机号验证
	if len(mobile) != 11 {
		return false
	}
	return true
}

// BindAndValidate 绑定并验证请求参数
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	// 根据 Content-Type 选择绑定器
	var err error
	method := c.Request.Method

	switch {
	case method == "GET":
		err = c.ShouldBindQuery(obj)
	case c.ContentType() == "multipart/form-data":
		err = c.ShouldBindWith(obj, binding.FormMultipart)
	default:
		err = c.ShouldBindJSON(obj)
	}

	// 处理验证错误
	if err != nil {
		processValidationError(c, err)
		return false
	}
	return true
}

// processValidationError 处理验证错误
func processValidationError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		// 翻译错误信息
		errs := make([]string, 0, len(validationErrs))
		for _, e := range validationErrs {
			errs = append(errs, e.Translate(trans))
		}
		ResponseError(c, CodeInvalidParams, strings.Join(errs, ", "))
		return
	}
	// 处理JSON解析错误
	ResponseError(c, CodeInvalidParams, fmt.Sprintf("参数解析错误: %s", err.Error()))
}

// ValidateQuery 验证查询参数
func ValidateQuery(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBindQuery(obj)
	if err != nil {
		processValidationError(c, err)
		return false
	}
	return true
}

// ValidateJSON 验证JSON参数
func ValidateJSON(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		processValidationError(c, err)
		return false
	}
	return true
}

// ValidateForm 验证表单参数
func ValidateForm(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBind(obj)
	if err != nil {
		processValidationError(c, err)
		return false
	}
	return true
}

// ValidateURI 验证路径参数
func ValidateURI(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBindUri(obj)
	if err != nil {
		processValidationError(c, err)
		return false
	}
	return true
}
