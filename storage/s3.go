package storage

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/nhost/hasura-storage/controller"
	"github.com/sirupsen/logrus"
)

func deptr[T any](p *T) T { //nolint:ireturn
	if p == nil {
		return *new(T)
	}
	return *p
}

type S3Error struct {
	Code    string `xml:"Code"`
	Message string `xml:"Message"`
}

func parseS3Error(resp *http.Response) *controller.APIError {
	var s3Error S3Error
	if err := xml.NewDecoder(resp.Body).Decode(&s3Error); err != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return controller.InternalServerError(
				fmt.Errorf("problem reading S3 error, status code %d: %w", resp.StatusCode, err),
			)
		}
		return controller.InternalServerError(
			fmt.Errorf( //nolint: goerr113
				"problem parsing S3 error, status code %d: %s",
				resp.StatusCode,
				b,
			),
		)
	}
	return controller.NewAPIError(
		resp.StatusCode,
		s3Error.Message,
		errors.New(s3Error.Message), //nolint: goerr113
		nil,
	)
}

type S3 struct {
	client     *s3.Client
	bucket     *string
	rootFolder string
	url        string
	logger     *logrus.Logger
}

func NewS3(
	client *s3.Client,
	bucket string,
	rootFolder string,
	url string,
	logger *logrus.Logger,
) *S3 {
	return &S3{
		client:     client,
		bucket:     aws.String(bucket),
		rootFolder: rootFolder,
		url:        url,
		logger:     logger,
	}
}

func (s *S3) PutFile(
	ctx context.Context,
	content io.ReadSeeker,
	filepath string,
	contentType string,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	// let's make sure we are in the beginning of the content
	if _, err := content.Seek(0, 0); err != nil {
		return "", controller.InternalServerError(
			fmt.Errorf("problem going to the beginning of the content: %w", err),
		)
	}

	object, err := s.client.PutObject(ctx,
		&s3.PutObjectInput{
			Body:        content,
			Bucket:      s.bucket,
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
	)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem putting object: %w", err))
	}

	return *object.ETag, nil
}

func (s *S3) GetFile(
	ctx context.Context,
	filepath string,
	headers http.Header,
) (*controller.File, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	object, err := s.client.GetObject(ctx,
		&s3.GetObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(key),
			// IfMatch:           new(string),
			// IfModifiedSince:   &time.Time{},
			// IfNoneMatch:       new(string),
			// IfUnmodifiedSince: &time.Time{},
			Range: aws.String(headers.Get("Range")),
		},
	)
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem getting object: %w", err))
	}

	status := http.StatusOK

	respHeaders := make(http.Header)
	if object.ContentRange != nil {
		respHeaders = http.Header{
			"Accept-Ranges": []string{"bytes"},
		}
		respHeaders["Content-Range"] = []string{*object.ContentRange}
		status = http.StatusPartialContent
	}

	return &controller.File{
		ContentType:   *object.ContentType,
		ContentLength: deptr(object.ContentLength),
		Etag:          *object.ETag,
		StatusCode:    status,
		Body:          object.Body,
		ExtraHeaders:  respHeaders,
	}, nil
}

func (s *S3) CreateGetObjectPresignedURL(
	ctx context.Context,
	filepath string,
	expire time.Duration,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	presignClient := s3.NewPresignClient(s.client)
	request, err := presignClient.PresignGetObject(ctx,
		&s3.GetObjectInput{ //nolint:exhaustivestruct
			Bucket: s.bucket,
			Key:    aws.String(key),
		},
		func(po *s3.PresignOptions) {
			po.Expires = expire
		},
	)
	if err != nil {
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	parts := strings.Split(request.URL, "?")
	if len(parts) != 2 { //nolint: mnd
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	return parts[1], nil
}

func (s *S3) GetFileWithPresignedURL(
	ctx context.Context, filepath, signature string, headers http.Header,
) (*controller.File, *controller.APIError) {
	if s.rootFolder != "" {
		filepath = s.rootFolder + "/" + filepath
	}
	url := fmt.Sprintf("%s/%s/%s?%s", s.url, *s.bucket, filepath, signature)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem creating request: %w", err))
	}
	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose //we are actually returning the body
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem getting file: %w", err))
	}

	if !(resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusPartialContent ||
		resp.StatusCode == http.StatusNotModified) {
		return nil, parseS3Error(resp)
	}

	respHeaders := make(http.Header)
	var length int64
	switch resp.StatusCode {
	case http.StatusOK, http.StatusPartialContent:
		respHeaders = http.Header{
			"Accept-Ranges": []string{"bytes"},
		}
		if resp.StatusCode == http.StatusPartialContent {
			respHeaders["Content-Range"] = []string{resp.Header.Get("Content-Range")}
		}

		length, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
		if err != nil {
			return nil, controller.InternalServerError(
				fmt.Errorf("problem parsing Content-Length: %w", err),
			)
		}
	}

	return &controller.File{
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: length,
		Etag:          resp.Header.Get("Etag"),
		StatusCode:    resp.StatusCode,
		Body:          resp.Body,
		ExtraHeaders:  respHeaders,
	}, nil
}

func (s *S3) PutFileWithPresignedURL(
	ctx context.Context, filepath, signature string, headers http.Header,
) (*httputil.ReverseProxy, *controller.APIError) {
	if s.rootFolder != "" {
		filepath = s.rootFolder + "/" + filepath
	}

	s3Url := fmt.Sprintf("%s/%s/%s?%s", s.url, *s.bucket, filepath, signature)

	u, err := url.Parse(s3Url)
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem parsing s3 url: %w", err))
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.Host = u.Host
		req.URL.Path = u.Path
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		ret := fmt.Sprintf("http proxy error %v", err)
		rw.Write([]byte(ret))
	}

	return proxy, nil
}

func (s *S3) DeleteFile(ctx context.Context, filepath string) *controller.APIError {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	if _, err := s.client.DeleteObject(ctx,
		&s3.DeleteObjectInput{
			Bucket: s.bucket,
			Key:    aws.String(key),
		}); err != nil {
		return controller.InternalServerError(fmt.Errorf("problem deleting file in s3: %w", err))
	}

	return nil
}

func (s *S3) ListFiles(ctx context.Context) ([]string, *controller.APIError) {
	objects, err := s.client.ListObjects(ctx,
		&s3.ListObjectsInput{
			Bucket: s.bucket,
			Prefix: aws.String(s.rootFolder + "/"),
		})
	if err != nil {
		return nil, controller.InternalServerError(
			fmt.Errorf("problem listing objects in s3: %w", err),
		)
	}

	res := make([]string, len(objects.Contents))
	for i, c := range objects.Contents {
		res[i] = strings.TrimPrefix(*c.Key, s.rootFolder+"/")
	}

	return res, nil
}

func (s *S3) CreateMultipartUpload(
	ctx context.Context,
	filepath string,
	contentType string,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	output, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      s.bucket,
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem create multipart upload in s3: %w", err))
	}

	return *output.UploadId, nil
}

func (s *S3) ListParts(
	ctx context.Context,
	filepath string,
	uploadId string,
) ([]controller.MultipartFragment, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	output, err := s.client.ListParts(ctx, &s3.ListPartsInput{
		Bucket:   s.bucket,
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	})
	if err != nil {
		return nil, controller.InternalServerError(fmt.Errorf("problem list parts in s3: %w", err))
	}

	parts := make([]controller.MultipartFragment, len(output.Parts))
	for _, part := range output.Parts {
		parts = append(parts, controller.MultipartFragment{
			ETag:         *part.ETag,
			LastModified: *part.LastModified,
			PartNumber:   *part.PartNumber,
			Size:         *part.Size,
		})
	}

	return parts, nil
}

func (s *S3) UploadPart(
	ctx context.Context,
	filepath string,
	uploadId string,
	partNumber int32,
	body io.ReadSeeker,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	output, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:     s.bucket,
		Key:        aws.String(key),
		UploadId:   aws.String(uploadId),
		PartNumber: aws.Int32(partNumber),
		Body:       body,
	})
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem upload part in s3: %w", err))
	}

	return *output.ETag, nil
}

func (s *S3) CreatePutObjectPresignedURL(
	ctx context.Context,
	filepath string,
	contentType string,
	expire time.Duration,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	presignClient := s3.NewPresignClient(s.client)
	request, err := presignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{ //nolint:exhaustivestruct
			Bucket:      s.bucket,
			Key:         aws.String(key),
			ContentType: aws.String(contentType),
		},
		func(po *s3.PresignOptions) {
			po.Expires = expire
		},
	)
	if err != nil {
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	logrus.Info(request.URL)

	parts := strings.Split(request.URL, "?")
	if len(parts) != 2 { //nolint: mnd
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	return parts[1], nil
}

func (s *S3) CreateUploadPartPresignedURL(
	ctx context.Context,
	filepath string,
	uploadId string,
	partNumber int32,
	expire time.Duration,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	presignClient := s3.NewPresignClient(s.client)
	request, err := presignClient.PresignUploadPart(ctx,
		&s3.UploadPartInput{ //nolint:exhaustivestruct
			Bucket:     s.bucket,
			Key:        aws.String(key),
			PartNumber: aws.Int32(partNumber),
			UploadId:   aws.String(uploadId),
		},
		func(po *s3.PresignOptions) {
			po.Expires = expire
		},
	)
	if err != nil {
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	parts := strings.Split(request.URL, "?")
	if len(parts) != 2 { //nolint: mnd
		return "", controller.InternalServerError(
			fmt.Errorf("problem generating pre-signed URL: %w", err),
		)
	}

	logrus.Info(request.URL)

	return parts[1], nil
}

func (s *S3) CompleteMultipartUpload(
	ctx context.Context,
	filepath string,
	uploadId string,
) (string, *controller.APIError) {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	listPartsOutput, listPartsErr := s.client.ListParts(ctx, &s3.ListPartsInput{
		Bucket:   s.bucket,
		Key:      aws.String(key),
		UploadId: aws.String(uploadId),
	})
	if listPartsErr != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem list parts in s3: %w", listPartsErr))
	}

	completedParts := make([]types.CompletedPart, 0, 10)
	for _, part := range listPartsOutput.Parts {
		if *part.Size == 0 {
			continue
		}
		completedParts = append(completedParts, types.CompletedPart{
			ChecksumCRC32:  part.ChecksumCRC32,
			ChecksumCRC32C: part.ChecksumCRC32C,
			ChecksumSHA1:   part.ChecksumSHA1,
			ChecksumSHA256: part.ChecksumSHA256,
			ETag:           part.ETag,
			PartNumber:     part.PartNumber,
		})
	}

	completeMultipartOutput, completeMultipartErr := s.client.CompleteMultipartUpload(ctx,
		&s3.CompleteMultipartUploadInput{
			Bucket:   s.bucket,
			Key:      aws.String(key),
			UploadId: aws.String(uploadId),
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: completedParts,
			},
		})
	if completeMultipartErr != nil {
		return "", controller.InternalServerError(fmt.Errorf("problem complete multipart upload in s3: %w", completeMultipartErr))
	}

	return *completeMultipartOutput.ETag, nil
}

func (s *S3) AbortMultipartUpload(
	ctx context.Context,
	filepath string,
	uploadId string,
) *controller.APIError {
	key, err := url.JoinPath(s.rootFolder, filepath)
	if err != nil {
		return controller.InternalServerError(fmt.Errorf("problem joining path: %w", err))
	}

	if _, err := s.client.AbortMultipartUpload(ctx,
		&s3.AbortMultipartUploadInput{
			Bucket:   s.bucket,
			Key:      aws.String(key),
			UploadId: aws.String(uploadId),
		}); err != nil {
		return controller.InternalServerError(fmt.Errorf("problem abort multipart upload in s3: %w", err))
	}

	return nil
}
