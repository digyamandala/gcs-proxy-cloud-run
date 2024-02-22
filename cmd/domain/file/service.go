package file

import (
	"context"

	uploaderclient "github.com/DomZippilli/gcs-proxy-cloud-function/backends/clients/uploader-client"
)

type Service interface {
	UploadFile(ctx context.Context, input FileUploadReq) (*FileUploadRes, error)
}
type service struct {
	uploaderClient uploaderclient.Client
}

func NewService(
	uploaderClient uploaderclient.Client,
) Service {
	return &service{
		uploaderClient: uploaderClient,
	}
}

func (ths *service) UploadFile(ctx context.Context, input FileUploadReq) (*FileUploadRes, error) {
	return nil, nil
}
