package filegateway

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/httpclient"
)

type fileGateway struct {
	client *httpclient.Client
}

func NewFileGateway(
	baseURL string,
	authServiceKey string,
) FileGateway {
	client := httpclient.NewClientBuilder().
		Context(context.Background()).
		BaseURL(baseURL).
		Header("X-Service-Key", authServiceKey).
		JSONContentType().
		RetryConfig(3, 1000).
		Build()

	return &fileGateway{
		client: client,
	}
}

func (f *fileGateway) DownloadFile(ctx context.Context, url string) (*DownloadResult, error) {
	if err := validateDownloadURL(ctx, url); err != nil {
		return nil, err
	}

	response, err := f.client.Post("/download").
		Context(ctx).
		Body(map[string]any{
			"url": url,
		}).
		DoStream()

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		err = response.BodyReader.Close()
		if err != nil {
			return nil, ErrDownloadFailed
		}
		return nil, ErrDownloadFailed
	}

	filename := response.Headers.Get("X-Filename")
	fileType := response.Headers.Get("X-Type")
	ext := response.Headers.Get("X-Ext")
	contentType := response.Headers.Get("Content-Type")
	contentLength := response.Headers.Get("Content-Length")

	var size int64
	if contentLength != "" {
		size = cast.ToInt64(contentLength)
	}

	return &DownloadResult{
		Filename:    filename,
		Type:        fileType,
		Ext:         ext,
		ContentType: contentType,
		Size:        size,
		DataReader:  response.BodyReader,
	}, nil
}

type uploadAPIResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data      struct {
		URL string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}

func (f *fileGateway) UploadFile(ctx context.Context, uploadID string, fileName string, reader io.Reader) (*UploadResult, error) {
	var builder strings.Builder
	builder.WriteString("/upload/")
	builder.WriteString(uploadID)
	uploadPath := builder.String()

	var successResponse uploadAPIResponse
	var errorResponse uploadAPIResponse
	response, err := f.client.Post(uploadPath).
		Context(ctx).
		MultipartFormBody().
		AddFileReader("file", fileName, reader).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, errors.New(errorResponse.Message)
	}

	return &UploadResult{
		URL: successResponse.Data.URL,
	}, nil
}
