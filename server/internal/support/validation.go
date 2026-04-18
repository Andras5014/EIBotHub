package support

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationMessage(err error, fieldLabels map[string]string) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return "请求参数格式不正确"
	}

	messages := make([]string, 0, len(validationErrors))
	for _, item := range validationErrors {
		field := item.Field()
		label := field
		if value, ok := fieldLabels[field]; ok {
			label = value
		}
		messages = append(messages, describeFieldError(label, item))
	}

	return strings.Join(messages, "；")
}

func describeFieldError(label string, item validator.FieldError) string {
	switch item.Tag() {
	case "required":
		return fmt.Sprintf("%s不能为空", label)
	case "email":
		return fmt.Sprintf("%s格式不正确", label)
	case "min":
		return fmt.Sprintf("%s至少%s个字符", label, item.Param())
	case "max":
		return fmt.Sprintf("%s不能超过%s个字符", label, item.Param())
	case "gte":
		return fmt.Sprintf("%s不能小于%s", label, item.Param())
	case "lte":
		return fmt.Sprintf("%s不能大于%s", label, item.Param())
	case "oneof":
		return fmt.Sprintf("%s取值不合法", label)
	default:
		return fmt.Sprintf("%s格式不正确", label)
	}
}
