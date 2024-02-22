package commonutils

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/apierror"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type key string

const (
	ZeroUUID             = "00000000-0000-0000-0000-000000000000"
	ReqIDKey         key = "Req-ID"
	XForwardedForKey key = "X-Forwarded-For"
	UserAgentKey     key = "User-Agent"
	CFRay            key = "CF-RAY"
	TrueClientIP     key = "True-Client-IP"
	LocaleID         key = "Locale-ID"
)

func ReqIDFromContext(ctx context.Context) string {
	id := ctx.Value(ReqIDKey)
	if id != nil {
		return id.(string)
	}
	return ZeroUUID
}

func HandleValidationError(err error) []apierror.APIError {
	var errorMap []apierror.APIError

	switch v := err.(type) {
	case *validator.InvalidValidationError:
		errorMap = append(errorMap, apierror.APIError{
			Code: apierror.CodeInvalidRequest,
			Desc: v.Error(),
		})
	case validator.ValidationErrors:
		for _, value := range v {
			var desc string
			switch value.Tag() {
			case "required_with":
				desc = fmt.Sprintf("Validation error: field %s is Required With %s", value.Field(), value.Param())
			case "oneof":
				desc = fmt.Sprintf("Validation error: field=%s, tag=%s, value=%s, Allowed values: %s", value.Field(), value.Tag(), value.Value(), value.Param())
			case "allow_character":
				params, _ := hex.DecodeString(strings.ReplaceAll(value.Param(), " ", ""))
				desc = fmt.Sprintf("Validation error: field=%s, tag=%s, value=%s, Allowed characters: alphanumeric and allowed additional characters: %s", value.Field(), value.Tag(), value.Value(), params)
			case "max":
				desc = fmt.Sprintf("Validation error: field %s has a max limit of %s", value.Field(), value.Param())
			default:
				desc = value.Error()
			}
			errorMap = append(errorMap, apierror.WithDesc(apierror.CodeInvalidRequest, desc))
		}
	default:
		errorMap = append(errorMap, apierror.APIError{
			Code: apierror.CodeInvalidRequest,
			Desc: v.Error(),
		})
	}

	return errorMap
}
func UserAgentFromContext(ctx context.Context) string {
	userAgent := ctx.Value(UserAgentKey)
	if userAgent != nil {
		return userAgent.(string)
	}
	return ""
}

func TrueClientIPFromContext(ctx context.Context) string {
	clientIP := ctx.Value(TrueClientIP)
	if clientIP != nil {
		return clientIP.(string)
	}
	return ""
}

func ContextWithReqID(ctx context.Context) context.Context {
	id, err := uuid.NewUUID()
	if err != nil {
		return context.WithValue(ctx, ReqIDKey, ZeroUUID)
	}
	return context.WithValue(ctx, ReqIDKey, id.String())
}
