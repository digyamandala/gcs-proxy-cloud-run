package uploaderclient

type APIModel[T interface{}] struct {
	Message string `json:"message"`
	Data    *T     `json:"data"`
}

type ErrorModel struct {
	ReqID   string `json:"reqId"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorAPIModel struct {
	Errors []ErrorModel `json:"errors"`
}

type RequestUploadSignedUrlReq struct {
	Identifier       string            `json:"identifier"`
	FileName         string            `json:"fileName"`
	ServiceName      string            `json:"serviceName"`
	IsPublic         bool              `json:"isPublic"`
	ImageMetadata    *ImageMetadata    `json:"imageMetadata"`
	AudioMetadata    *AudioMetadata    `json:"audioMetadata"`
	DocumentMetadata *DocumentMetadata `json:"documentMetadata"`
	VideoMetadata    *VideoMetadata    `json:"videoMetadata"`
	Project          Project           `json:"project"`
}

type UploadStatusReq struct {
	Tokens []string `json:"tokens"`
}
type RequestUploadSignedUrlRes struct {
	Identifier string `json:"identifier"`
	SignedUrl  string `json:"signedUrl"`
	JWTToken   string `json:"jwtToken"`
	Expiry     int    `json:"expiry"`
}

type RequestDownloadUrlReq struct {
	Token          []string `json:"token"`
	ExpiryInSecond int      `json:"expiryInSecond"`
}

type RequestDownloadUrlRes struct {
	Token     string `json:"token"`
	SignedUrl string `json:"signedUrl"`
	Expiry    int    `json:"expiry"`
	PublicUrl string `json:"publicUrl"`
}

type CheckUploadStatusWithWaitReq struct {
	Tokens []string `json:"tokens"`
}

type ImageMetadata struct {
	MinSize     int64  `json:"minSize"`
	MaxSize     int64  `json:"maxSize"`
	MinWidth    int64  `json:"minWidth"`
	MaxWidth    int64  `json:"maxWidth"`
	MinHeight   int64  `json:"minHeight"`
	MaxHeight   int64  `json:"maxHeight"`
	ContentType string `json:"contentType"`
}

type AudioMetadata struct {
	MinSize     int64  `json:"minSize"`
	MaxSize     int64  `json:"maxSize"`
	Duration    int64  `json:"duration"`
	ContentType string `json:"contentType"`
}

type DocumentMetadata struct {
	MinSize     int64  `json:"minSize"`
	MaxSize     int64  `json:"maxSize"`
	ContentType string `json:"contentType"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

type VideoMetadata struct {
	MinSize     int64  `json:"minSize"`
	MaxSize     int64  `json:"maxSize"`
	Duration    int64  `json:"duration"`
	ContentType string `json:"contentType"`
}

type GetUploadStatusInput struct {
	Token string `json:"token"`
}

type Status string

const (
	INITIATE_UPLOAD               Status = "INITIATE_UPLOAD"
	UPLOAD_ON_METADATA_VALIDATION Status = "UPLOAD_ON_METADATA_VALIDATION"
	UPLOAD_ON_MALWARE_SCANNING    Status = "UPLOAD_ON_MALWARE_SCANNING"
	UPLOAD_VALIDATION_FAILED      Status = "UPLOAD_VALIDATION_FAILED"
	UPLOAD_VIRUS_DETECTED         Status = "UPLOAD_VIRUS_DETECTED"
	UPLOAD_SUCCESS                Status = "UPLOAD_SUCCESS"
)
