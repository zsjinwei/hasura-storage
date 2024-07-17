package controller

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// this type is used to ensure we respond consistently no matter the case.
type CreateFileMultipartUploadResponse struct {
	Files []FileMetadata `json:"files,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type fileMultipartUploadData struct {
	Size        int64  `json:"size"`
	ChunkSize   int64  `json:"chunkSize"`
	ChunkCount  int64  `json:"chunkCount"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

type createFileMultipartUploadRequest struct {
	BucketID     string                    `json:"bucketID"`
	ObjectPrefix string                    `json:"objectPrefix"`
	Files        []fileMultipartUploadData `json:"files"`
}

func parseCreateFileMultipartUploadRequest(ctx *gin.Context) (createFileMultipartUploadRequest, *APIError) {
	json := createFileMultipartUploadRequest{}

	err := ctx.BindJSON(&json)
	if err != nil {
		return createFileMultipartUploadRequest{}, BadDataError(errors.New("request data is invalid"), "request data is invalid")
	}

	if json.BucketID == "" {
		json.BucketID = "default"
	}

	if len(json.Files) == 0 {
		errMsg := "file parameter file len must be greater than zero"
		return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
	}

	for i, file := range json.Files {
		if file.Size <= 0 {
			errMsg := fmt.Sprintf("file %d parameter size must be greater than zero", i)
			return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
		}

		if file.ChunkSize <= 0 {
			errMsg := fmt.Sprintf("file %d parameter chunkSize must be greater than zero", i)
			return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
		}

		if file.Size < file.ChunkSize {
			errMsg := fmt.Sprintf("file %d size must be greater than chunkSize", i)
			return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
		}

		if file.FileName == "" {
			errMsg := fmt.Sprintf("file %d parameter fileName must be not empty", i)
			return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
		}

		if file.ContentType == "" {
			json.Files[i].ContentType = "application/octet-stream"
		}

		json.Files[i].ChunkCount = int64(math.Ceil(float64(file.Size) / float64(file.ChunkSize)))

		if json.Files[i].ChunkCount > 1 && file.ChunkSize < 5*1024*1024 {
			errMsg := fmt.Sprintf("file %d parameter chunkSize must be greater than 5MB", i)
			return createFileMultipartUploadRequest{}, BadDataError(errors.New(errMsg), errMsg)
		}
	}

	return json, nil
}

func (ctrl *Controller) createFileMultipartUploadProcess(ctx *gin.Context) ([]FileMetadata, *APIError) {
	req, apiErr := parseCreateFileMultipartUploadRequest(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	bucket, err := ctrl.metadataStorage.GetBucketByID(
		ctx,
		req.BucketID,
		http.Header{"x-hasura-admin-secret": []string{ctrl.hasuraAdminSecret}},
	)
	if err != nil {
		return nil, err
	}

	minSize := bucket.MinUploadFile
	maxSize := bucket.MaxUploadFile

	fileMetadatas := make([]FileMetadata, 0, len(req.Files))

	for _, file := range req.Files {
		if minSize > int(file.Size) {
			return nil, FileTooSmallError(file.FileName, int(file.Size), minSize)
		} else if int(file.Size) > maxSize {
			return nil, FileTooBigError(file.FileName, int(file.Size), maxSize)
		}

		fileId := uuid.New().String()
		objectKey, joinErr := url.JoinPath(req.ObjectPrefix, fileId)
		if joinErr != nil {
			return nil, InternalServerError(
				fmt.Errorf("problem joining path: %w", joinErr),
			)
		}

		uploadId, cErr := ctrl.contentStorage.CreateMultipartUpload(ctx, objectKey, file.ContentType)
		if cErr != nil {
			return nil, cErr
		}

		if err := ctrl.metadataStorage.InitializeFile(
			ctx, fileId, file.FileName, file.Size, bucket.ID, file.ContentType, objectKey, file.ChunkSize, file.ChunkCount, uploadId, ctx.Request.Header,
		); err != nil {
			return nil, err
		}

		fileMetadata, apiErr := ctrl.metadataStorage.GetFileByID(ctx, fileId, http.Header{"x-hasura-admin-secret": []string{ctrl.hasuraAdminSecret}})
		if apiErr != nil {
			return nil, apiErr
		}

		fileMetadatas = append(fileMetadatas, fileMetadata)
	}

	return fileMetadatas, nil
}

func (ctrl *Controller) CreateFileMultipartUpload(ctx *gin.Context) {
	filesMetadatas, apiErr := ctrl.createFileMultipartUploadProcess(ctx)
	if apiErr != nil {
		_ = ctx.Error(fmt.Errorf("problem processing request: %w", apiErr))

		ctx.JSON(apiErr.statusCode, CreateFileMultipartUploadResponse{filesMetadatas, apiErr.PublicResponse()})

		return
	}

	ctx.JSON(http.StatusOK, CreateFileMultipartUploadResponse{filesMetadatas, nil})
}
