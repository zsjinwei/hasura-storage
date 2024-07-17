package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// this type is used to ensure we respond consistently no matter the case.
type CompleteFileMultipartUploadResponse struct {
	Error *ErrorResponse `json:"error,omitempty"`
	File  FileMetadata   `json:"file"`
}

type completeFileMultipartUploadRequest struct {
	FileID string
}

func parseCompleteFileMultipartUploadRequest(ctx *gin.Context) (completeFileMultipartUploadRequest, *APIError) {
	return completeFileMultipartUploadRequest{
		FileID: ctx.Param("id"),
	}, nil
}

func checkFileMultipartUploadParts(fileMetadata FileMetadata, parts []MultipartFragment) (bool, *APIError) {
	completedChunkCount := int64(0)
	completedChunkSize := int64(0)

	for _, part := range parts {
		if part.Size == 0 {
			continue
		}
		completedChunkCount += 1
		completedChunkSize += part.Size
	}

	if fileMetadata.ChunkCount != completedChunkCount {
		errMsg := fmt.Sprintf("multipart upload is not completes: chunk count not matched(%d != %d)", fileMetadata.ChunkCount, completedChunkCount)
		return false, BadDataError(errors.New(errMsg), errMsg)
	}

	if fileMetadata.Size != completedChunkSize {
		errMsg := fmt.Sprintf("multipart upload is not completes: chunk size not matched(%d != %d)", fileMetadata.Size, completedChunkSize)
		return false, BadDataError(errors.New(errMsg), errMsg)
	}

	return true, nil
}

func (ctrl *Controller) completeFileMultipartUploadProcess(ctx *gin.Context) (FileMetadata, *APIError) {
	req, apiErr := parseCompleteFileMultipartUploadRequest(ctx)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	fileMetadata, _, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(), req.FileID, false, ctx.Request.Header,
	)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	if fileMetadata.UploadID == "" {
		return FileMetadata{},
			apiErr.ExtendError(
				"no multipart uploadId for file " + fileMetadata.Name,
			)
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	parts, err := ctrl.contentStorage.ListParts(ctx, objectKey, fileMetadata.UploadID)
	if err != nil {
		return FileMetadata{}, err
	}

	_, checkErr := checkFileMultipartUploadParts(fileMetadata, parts)
	if checkErr != nil {
		return FileMetadata{}, checkErr
	}

	etag, apiErr := ctrl.contentStorage.CompleteMultipartUpload(ctx, objectKey, fileMetadata.UploadID)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	metadata, apiErr := ctrl.metadataStorage.PopulateMetadata(
		ctx,
		fileMetadata.ID, fileMetadata.Name, fileMetadata.Size, fileMetadata.BucketID, etag, true, fileMetadata.MimeType, objectKey, fileMetadata.ChunkSize, fileMetadata.ChunkCount, fileMetadata.UploadID, fileMetadata.Metadata,
		http.Header{"x-hasura-admin-secret": []string{ctrl.hasuraAdminSecret}},
	)
	if apiErr != nil {
		return FileMetadata{}, apiErr.ExtendError(
			"problem populating file metadata for file " + fileMetadata.Name,
		)
	}

	return metadata, nil
}

func (ctrl *Controller) CompleteFileMultipartUpload(ctx *gin.Context) {
	fileMetadata, apiErr := ctrl.completeFileMultipartUploadProcess(ctx)
	if apiErr != nil {
		_ = ctx.Error(fmt.Errorf("problem processing request: %w", apiErr))

		ctx.JSON(apiErr.statusCode, CompleteFileMultipartUploadResponse{apiErr.PublicResponse(), FileMetadata{}})

		return
	}

	ctx.JSON(http.StatusOK, CompleteFileMultipartUploadResponse{nil, fileMetadata})
}
