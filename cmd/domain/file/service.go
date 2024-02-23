package file

import (
	"context"
	"errors"

	uploaderclient "github.com/DomZippilli/gcs-proxy-cloud-function/backends/clients/uploader-client"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/commonutils"
	"github.com/ztrue/tracerr"
)

type Service interface {
	UploadFile(ctx context.Context, input FileUploadReq) (string, error)
	DownloadFile(ctx context.Context, input string) (*uploaderclient.RequestDownloadUrlRes, error)
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

func (ths *service) UploadFile(ctx context.Context, input FileUploadReq) (string, error) {
	//UNCOMMENT FOR TESTING FILE
	// file, _ := os.Open("/Users/keziaaurelia/Downloads/download.jpeg")
	// defer file.Close()

	// // Get the file size
	// stat, _ := file.Stat()

	// // Read the file into a byte slice
	// bs := make([]byte, stat.Size())
	// _, _ = bufio.NewReader(file).Read(bs)
	mapUploadSignedUrlReq := make(map[string][]uploaderclient.RequestUploadSignedUrlReq)
	for _, req := range input.UploadSignedUrlReq {
		requestUploadSignedUrlReq := uploaderclient.RequestUploadSignedUrlReq{
			Identifier: req.Identifier,
			FileName:   req.FileName,
			IsPublic:   req.IsPublic,
		}

		switch input.Type {
		case IMAGE:
			imageMetadata := VALIDATION_IMAGE_METADATA[input.Type]
			imageMetadata.ContentType = req.ContentType
			requestUploadSignedUrlReq.ImageMetadata = &imageMetadata
		case VIDEO:
			videoMetadata := VALIDATION_VIDEO_METADATA[input.Type]
			videoMetadata.ContentType = req.ContentType
			requestUploadSignedUrlReq.VideoMetadata = &videoMetadata
		case DOCUMENT:
			documentMetadata := VALIDATION_DOCUMENT_METADATA[input.Type]
			documentMetadata.ContentType = req.ContentType
			requestUploadSignedUrlReq.DocumentMetadata = &documentMetadata
		case BULK_ACTION:
			if req.ContentType != XLSX_CONTENT_TYPE {
				return "", errors.New("invalid content type")
			}
			documentMetadata := VALIDATION_DOCUMENT_METADATA[input.Type]
			documentMetadata.ContentType = req.ContentType
			requestUploadSignedUrlReq.DocumentMetadata = &documentMetadata
		}
		mapUploadSignedUrlReq[req.ContentType] = append(mapUploadSignedUrlReq[req.ContentType], requestUploadSignedUrlReq)
	}
	uploadURL := []uploaderclient.RequestUploadSignedUrlRes{}
	for _, req := range mapUploadSignedUrlReq {
		resp, err := ths.uploaderClient.RequestUploadSignedUrl(
			commonutils.ReqIDFromContext(ctx),
			req,
		)
		if err != nil {
			return "", tracerr.Wrap(err)
		}
		uploadURL = append(uploadURL, resp...)
	}

	fileID, err := ths.uploaderClient.VerifyAndDecodeToken(uploadURL[0].JWTToken)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	//UNCOMMENT TO TESTING FILE
	// err = ths.uploaderClient.UploadFile(commonutils.ReqIDFromContext(ctx), bs, uploadURL[0].SignedUrl)
	err = ths.uploaderClient.UploadFile(commonutils.ReqIDFromContext(ctx), input.UploadSignedUrlReq[0].DocumentByte, uploadURL[0].SignedUrl)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	return fileID, nil
}

func (ths *service) DownloadFile(ctx context.Context, input string) (*uploaderclient.RequestDownloadUrlRes, error) {
	tmp := []string{}
	tmp = append(tmp, input)
	file, _ := ths.uploaderClient.RequestDownloadUrl(commonutils.ReqIDFromContext(ctx), uploaderclient.RequestDownloadUrlReq{
		Token:          tmp,
		ExpiryInSecond: 10000,
	})
	return &file[0], nil
}
