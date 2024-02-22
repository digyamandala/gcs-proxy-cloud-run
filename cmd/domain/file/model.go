package file

import uploaderclient "github.com/DomZippilli/gcs-proxy-cloud-function/backends/clients/uploader-client"

const (
	BULK_ACTION       = "BULK_ACTION"
	IMAGE             = "IMAGE"
	VIDEO             = "VIDEO"
	DOCUMENT          = "DOCUMENT"
	XLSX_CONTENT_TYPE = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

type FileUploadReq struct {
	Type               string               `json:"type" validate:"required,oneof=BULK_ACTION IMAGE VIDEO DOCUMENT"`
	UploadSignedUrlReq []UploadSignedUrlReq `json:"uploadSignedUrlReq" validation:"dive"`
}

var VALIDATION_IMAGE_METADATA = map[string]uploaderclient.ImageMetadata{
	IMAGE: {MinSize: 1, MaxSize: 40000000, MinWidth: 300, MaxWidth: 2048, MinHeight: 300, MaxHeight: 2048},
}

var VALIDATION_VIDEO_METADATA = map[string]uploaderclient.VideoMetadata{
	VIDEO: {MinSize: 1, MaxSize: 50000000, Duration: 120},
}

var VALIDATION_DOCUMENT_METADATA = map[string]uploaderclient.DocumentMetadata{
	DOCUMENT:    {MinSize: 1, MaxSize: 2000000},
	BULK_ACTION: {MinSize: 1, MaxSize: 50000000},
}

type UploadSignedUrlReq struct {
	ContentType  string `json:"contentType" validate:"required,oneof=application/vnd.openxmlformats-officedocument.spreadsheetml.sheet image/jpg image/jpeg image/png video/mp4 video/mov application/pdf"`
	Identifier   string `json:"identifier"`
	FileName     string `json:"fileName"`
	IsPublic     bool   `json:"isPublic"`
	DocumentByte []byte `json:"documentByte"`
}

type VerifyAndDecodeTokenReq struct {
	Token string
}

type RequestDownloadUrlReq struct {
	Token    []string
	IsPublic bool
}
