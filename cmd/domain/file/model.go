package file

import uploaderclient "github.com/DomZippilli/gcs-proxy-cloud-function/backends/clients/uploader-client"

const (
	BULK_ACTION       = "BULK_ACTION"
	IMAGE             = "IMAGE"
	VIDEO             = "VIDEO"
	DOCUMENT          = "DOCUMENT"
	XLSX_CONTENT_TYPE = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	SHEET             = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	IMAGE_JPG         = "image/jpg"
	IMAGE_JPEG        = "image/jpeg"
	IMAGE_PNG         = "image/png"
	VIDEO_MP4         = "video/mp4"
	VIDEO_MOV         = "video/mov"
	DOCUMENT_PDF      = "application/pdf"
	DOCUMENT_RHS      = "application/octet-stream"
)

type FileUploadReq struct {
	UploadSignedUrlReq []UploadSignedUrlReq `json:"uploadSignedUrlReq" validation:"dive"`
}

var VALIDATION_IMAGE_METADATA = map[string]uploaderclient.ImageMetadata{
	IMAGE: {MinSize: 1, MaxSize: 40000000, MinWidth: 1, MaxWidth: 2048, MinHeight: 1, MaxHeight: 2048},
}

var VALIDATION_VIDEO_METADATA = map[string]uploaderclient.VideoMetadata{
	VIDEO: {MinSize: 1, MaxSize: 50000000, Duration: 120},
}

var VALIDATION_DOCUMENT_METADATA = map[string]uploaderclient.DocumentMetadata{
	DOCUMENT:    {MinSize: 1, MaxSize: 2000000},
	BULK_ACTION: {MinSize: 1, MaxSize: 50000000},
}

type UploadSignedUrlReq struct {
	ContentType string `json:"contentType" validate:"required,oneof=application/vnd.openxmlformats-officedocument.spreadsheetml.sheet image/jpg image/jpeg image/png video/mp4 video/mov application/pdf"`
	Identifier  string `json:"identifier"`
	FileName    string `json:"fileName"`
	IsPublic    bool   `json:"isPublic"`
}
type UploadSignedUrlRes struct {
	SignedURL string `json:"signedUrl"`
	FileID    string `json:"fileId"`
}

type VerifyAndDecodeTokenReq struct {
	Token string `json:"token"`
}
type UploadStatusReq struct {
	Tokens []string `json:"tokens"`
}

type RequestDownloadUrlReq struct {
	Token    []string
	IsPublic bool
}
