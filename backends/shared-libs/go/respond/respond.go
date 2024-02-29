package respond

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/apierror"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/commonutils"
)

type APIModel[T interface{}] struct {
	Data   T               `json:"data"`
	Errors []ErrorAPIModel `json:"errors"`
}

type ErrorAPIModel struct {
	ReqID   string `json:"reqId"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	model := APIModel[interface{}]{Data: data}
	js, err := json.Marshal(model)
	if err != nil {
		Error(w, context.Background(), apierror.WithDesc(apierror.CodeInternalServerError, "Unknown error"), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(js); err != nil {
		log.Fatalf("%v", err)
	}
}

func Error(w http.ResponseWriter, ctx context.Context, errModel apierror.APIError, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	model := APIModel[interface{}]{Data: nil, Errors: []ErrorAPIModel{ErrorAPIModel{
		ReqID:   commonutils.ReqIDFromContext(ctx),
		Code:    errModel.Code,
		Message: errModel.Desc,
	}}}
	js, err := json.Marshal(model)
	if err != nil {
		Error(w, ctx, apierror.WithDesc(apierror.CodeInternalServerError, "Unknown error"), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(js); err != nil {
		log.Fatalf("%v", err)
	}
}

func MultiError(w http.ResponseWriter, ctx context.Context, apiErrors []apierror.APIError, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errors []ErrorAPIModel
	reqId := commonutils.ReqIDFromContext(ctx)
	for _, apiError := range apiErrors {
		errors = append(errors, ErrorAPIModel{
			ReqID:   reqId,
			Code:    apiError.Code,
			Message: apiError.Desc,
		})
	}

	model := APIModel[interface{}]{Data: nil, Errors: errors}
	js, err := json.Marshal(model)
	if err != nil {
		Error(w, ctx, apierror.WithDesc(apierror.CodeInternalServerError, "Unknown error"), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(js); err != nil {
		log.Fatalf("%v", err)
	}
}
