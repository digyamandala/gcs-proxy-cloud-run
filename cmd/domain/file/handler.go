package file

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/apierror"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/logger"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/respond"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	UploadFile(w http.ResponseWriter, req *http.Request)
	DownloadFile(w http.ResponseWriter, req *http.Request)
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{
		svc: svc,
	}
}

func (ths *handler) UploadFile(w http.ResponseWriter, req *http.Request) {
	var input FileUploadReq
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	res, err := ths.svc.UploadFile(req.Context(), input)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeEntityNotFound, "Category or Product not found"), http.StatusNotFound)
			return
		}
		return
	}
	respond.Success(w, res, http.StatusOK)
}

func (ths *handler) DownloadFile(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	res, err := ths.svc.DownloadFile(req.Context(), id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeEntityNotFound, "Category or Product not found"), http.StatusNotFound)
			return
		}
		return
	}
	respond.Success(w, res, http.StatusOK)
}
