package file

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/apierror"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/logger"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/respond"
	"github.com/DomZippilli/gcs-proxy-cloud-function/common"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler interface {
	HealthCheck(w http.ResponseWriter, req *http.Request)
	UploadFile(w http.ResponseWriter, req *http.Request)
	DownloadFile(w http.ResponseWriter, req *http.Request)
	VerifyAndDecodeToken(w http.ResponseWriter, req *http.Request)
	UploadStatus(w http.ResponseWriter, req *http.Request)
	HandlingOption(w http.ResponseWriter, req *http.Request)
}

type handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return &handler{
		svc: svc,
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func (ths *handler) HealthCheck(w http.ResponseWriter, req *http.Request) {
	respond.Success(w, "HEALTHY", http.StatusOK)
}

func (ths *handler) UploadFile(w http.ResponseWriter, req *http.Request) {
	if stringUtils.IsEmpty(req.Header.Get("x-lpse-id")) {
		http.Error(w, "x-lpse-id header not found", http.StatusBadRequest)
		return
	}
	enableCors(&w)
	var input FileUploadReq
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	lpseId := req.Header.Get("x-lpse-id")
	for i := 0; i < len(input.UploadSignedUrlReq); i++ {
		normalizedPath := common.NormalizePath(lpseId, input.UploadSignedUrlReq[i].FileName)
		log.Info().Msgf("normalized path %s", normalizedPath)
		input.UploadSignedUrlReq[i].FileName = normalizedPath
	}

	res, err := ths.svc.UploadFile(req.Context(), input)

	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInternalServerError, "Internal Server Error"), http.StatusBadRequest)
		return
	}
	respond.Success(w, res, http.StatusOK)
}

func (ths *handler) VerifyAndDecodeToken(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	var input VerifyAndDecodeTokenReq
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	res, err := ths.svc.VerifyAndDecodeToken(req.Context(), input)

	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInternalServerError, "Internal Server Error"), http.StatusBadRequest)
		return
	}
	respond.Success(w, res, http.StatusOK)
}

func (ths *handler) DownloadFile(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	id := chi.URLParam(req, "id")
	if id == "" {
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	res, err := ths.svc.DownloadFile(req.Context(), id)

	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInternalServerError, "Internal Server Error"), http.StatusBadRequest)
		return
	}
	if strings.Compare(res.SignedUrl, "") != 0 {
		log.Info().Msgf("Download fileID %s will be redirected to signedUrl %s", id, res.SignedUrl)
		http.Redirect(w, req, res.SignedUrl, http.StatusMovedPermanently)
	} else {
		log.Info().Msgf("Download fileID %s will be redirected to publicUrl %s", id, res.PublicUrl)
		http.Redirect(w, req, res.PublicUrl, http.StatusMovedPermanently)
	}
	respond.Success(w, res, http.StatusOK)
}

func (ths *handler) UploadStatus(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	var input UploadStatusReq
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.WithDesc(apierror.CodeInvalidRequest, "Invalid request"), http.StatusBadRequest)
		return
	}
	err = ths.svc.UploadStatus(req.Context(), input)

	if err != nil {
		logger.Warn(req.Context(), "%v", err)
		respond.Error(w, req.Context(), apierror.FromError(err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	respond.Success(w, nil, http.StatusOK)
}

func (ths *handler) HandlingOption(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	respond.Success(w, true, http.StatusOK)
}
