package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/zayn1510/goarchi/app/resources"
)

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	// Inisialisasi translator
	eng := en.New()
	uni := ut.New(eng, eng)

	trans, _ := uni.GetTranslator("en")
	validate = validator.New()

	enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return name
	})

	translator = trans
}

func Validate(data any) (error, map[string]string) {
	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Translate(translator)
	}

	return err, errors
}

func HandleBindError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, io.EOF):
		resources.BadRequest(ctx, "Body JSON is empty")
		return

	case strings.Contains(err.Error(), "invalid character"):
		resources.BadRequest(ctx, "Invalid JSON format")
		return

	default:
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			resources.BadRequest(ctx, fmt.Sprintf("Syntax error at byte offset %d", syntaxErr.Offset))
		case errors.As(err, &typeErr):
			resources.BadRequest(ctx, fmt.Sprintf(
				"Invalid type for field '%s': expected %s but got %s",
				typeErr.Field, typeErr.Type.String(), typeErr.Value,
			))
		default:
			resources.BadRequest(ctx, err.Error())
		}
	}
}
