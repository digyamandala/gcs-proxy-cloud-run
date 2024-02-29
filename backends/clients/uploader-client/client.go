package uploaderclient

//go:generate mockgen -source=./client.go -destination=../../mocks/uploader/client_mock.go

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/apierror"
	"github.com/DomZippilli/gcs-proxy-cloud-function/backends/shared-libs/go/logger"
	"github.com/go-resty/resty/v2"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ztrue/tracerr"
)

type Client interface {
	RequestUploadSignedUrl(reqID string, req []RequestUploadSignedUrlReq) ([]RequestUploadSignedUrlRes, error)
	RequestDownloadUrl(reqID string, req RequestDownloadUrlReq) ([]RequestDownloadUrlRes, error)
	VerifyAndDecodeToken(token string) (string, error)
	RequestDownloadUrlWithWait(ctx context.Context, reqID string, req RequestDownloadUrlReq) ([]RequestDownloadUrlRes, error)
	UploadFile(reqID string, file interface{}, signedUrl string) error
}

type client struct {
	baseURL      string
	publicJwkSet jwk.Set
	restyClient  *resty.Client
}

func NewClient(baseURL string, publicKey []byte) (*client, error) {
	restyClient := resty.New()
	restyClient.OnAfterResponse(func(rc *resty.Client, resp *resty.Response) error {
		if resp.IsError() {
			req := resp.Request

			logger.Error(req.Context(), "req: %#v; res: %#v", req, resp.String())
		}

		return nil
	})

	newClient := &client{
		baseURL:     baseURL,
		restyClient: restyClient,
	}

	if publicKey == nil {
		jwkSet, err := newClient.getJwkSet()
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		newClient.publicJwkSet = jwkSet
	}

	return newClient, nil
}

func (c *client) RequestUploadSignedUrl(reqID string, req []RequestUploadSignedUrlReq) (res []RequestUploadSignedUrlRes, err error) {
	url := fmt.Sprintf("%s/upload/bulkRequest", c.baseURL)
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&APIModel[[]RequestUploadSignedUrlRes]{}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Request-ID", reqID).
		Post(url)
	if err != nil {
		err = fmt.Errorf("error when doing api call for uploader request upload signed url: %w", err)
		return
	}

	if resp.IsError() {
		var errorModel ErrorAPIModel
		err = json.Unmarshal(resp.Body(), &errorModel)
		if err != nil {
			return
		}

		err = apierror.WithDesc(errorModel.Errors[0].Code, errorModel.Errors[0].Message)
		return
	}

	var model APIModel[[]RequestUploadSignedUrlRes]
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		return
	}

	return *model.Data, nil
}

func (c *client) RequestDownloadUrl(reqID string, req RequestDownloadUrlReq) (res []RequestDownloadUrlRes, err error) {
	url := fmt.Sprintf("%s/download/request", c.baseURL)
	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&APIModel[[]RequestDownloadUrlRes]{}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Request-ID", reqID).
		Post(url)
	if err != nil {
		err = fmt.Errorf("error when doing api call for uploader request download signed url: %w", err)
		return
	}

	if resp.IsError() {
		var errorModel ErrorAPIModel
		err = json.Unmarshal(resp.Body(), &errorModel)
		if err != nil {
			return
		}

		err = apierror.WithDesc(errorModel.Errors[0].Code, errorModel.Errors[0].Message)
		return
	}

	var model APIModel[[]RequestDownloadUrlRes]
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		return
	}

	return *model.Data, nil
}

func (c *client) VerifyAndDecodeToken(token string) (string, error) {
	var jwtToken jwt.Token
	var err error
	if c.publicJwkSet != nil {
		jwtToken, err = jwt.ParseString(token, jwt.WithKeySet(c.publicJwkSet), jwt.WithValidate(true))
		if err != nil {
			return "", err
		}
	} else {
		jwtToken, err = jwt.ParseString(token)
		if err != nil {
			return "", err
		}
	}

	uploadToken, isExist := jwtToken.Get("token")
	if !isExist {
		return "", fmt.Errorf("upload token not exist in jwt token")
	}

	return fmt.Sprintf("%v", uploadToken), nil
}

func (c *client) RequestDownloadUrlWithWait(ctx context.Context, reqID string, req RequestDownloadUrlReq) (res []RequestDownloadUrlRes, err error) {
	urlCheckUpload := fmt.Sprintf("%s/upload/check", c.baseURL)
	url := fmt.Sprintf("%s/download/request", c.baseURL)

	checkUploadReq := CheckUploadStatusWithWaitReq{Tokens: req.Token}
	checkUploadRes, err := c.restyClient.R().
		SetBody(checkUploadReq).
		SetHeader("Content-Type", "application/json").
		SetHeader("Request-ID", reqID).
		Post(urlCheckUpload)
	if err != nil {
		err = apierror.FromError(fmt.Errorf("error when doing api call for /upload/check: %w", err))
		return
	}

	if checkUploadRes.IsError() {
		var errorModel ErrorAPIModel
		err = json.Unmarshal(checkUploadRes.Body(), &errorModel)
		if err != nil {
			return
		}

		err = apierror.WithDesc(errorModel.Errors[0].Code, errorModel.Errors[0].Message)
		return
	}

	resp, err := c.restyClient.R().
		SetBody(req).
		SetResult(&APIModel[[]RequestDownloadUrlRes]{}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Request-ID", reqID).
		Post(url)
	if err != nil {
		err = fmt.Errorf("error when doing api call for uploader request download signed url: %w", err)
		return
	}

	if resp.IsError() {
		var errorModel ErrorAPIModel
		err = json.Unmarshal(resp.Body(), &errorModel)
		if err != nil {
			return
		}

		err = apierror.WithDesc(errorModel.Errors[0].Code, errorModel.Errors[0].Message)
		return
	}

	var model APIModel[[]RequestDownloadUrlRes]
	err = json.Unmarshal(resp.Body(), &model)
	if err != nil {
		return
	}

	return *model.Data, nil
}

func (c *client) UploadFile(reqID string, file interface{}, signedUrl string) error {
	_, err := c.restyClient.R().
		SetBody(file).
		SetHeader("Content-Range", "bytes 0-*/*").
		Put(signedUrl)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getJwkSet() (jwk.Set, error) {
	url := fmt.Sprintf("%s/.well-known/jwks", c.baseURL)
	resp, err := c.restyClient.R().
		Get(url)
	if err != nil {
		err = fmt.Errorf("error get public jwk set: %w", err)
		return nil, tracerr.Wrap(err)
	}

	if resp.IsError() {
		var errorModel ErrorAPIModel
		err = json.Unmarshal(resp.Body(), &errorModel)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		err = apierror.WithDesc(errorModel.Errors[0].Code, errorModel.Errors[0].Message)
		return nil, tracerr.Wrap(err)
	}
	jwkSet, err := jwk.Parse(resp.Body())
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return jwkSet, nil
}
