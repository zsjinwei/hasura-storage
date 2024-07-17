package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// this type is used to ensure we respond consistently no matter the case.
type AbortFileMultipartUploadResponse struct {
	Error *ErrorResponse `json:"error,omitempty"`
}

type abortFileMultipartUploadRequest struct {
	FileID string
}

func parseAbortFileMultipartUploadRequest(ctx *gin.Context) (abortFileMultipartUploadRequest, *APIError) {
	return abortFileMultipartUploadRequest{
		FileID: ctx.Param("id"),
	}, nil
}

func (ctrl *Controller) abortFileMultipartUploadProcess(ctx *gin.Context) *APIError {
	req, apiErr := parseAbortFileMultipartUploadRequest(ctx)
	if apiErr != nil {
		return apiErr
	}

	fileMetadata, _, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(), req.FileID, false, ctx.Request.Header,
	)
	if apiErr != nil {
		return apiErr
	}

	if fileMetadata.UploadID == "" {
		return apiErr.ExtendError(
			"no multipart uploadId for file " + fileMetadata.Name,
		)
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	abortErr := ctrl.contentStorage.AbortMultipartUpload(ctx, objectKey, fileMetadata.UploadID)
	if abortErr != nil {
		return abortErr
	}

	delErr := ctrl.metadataStorage.DeleteFileByID(
		ctx,
		fileMetadata.ID,
		http.Header{"x-hasura-admin-secret": []string{ctrl.hasuraAdminSecret}},
	)
	if delErr != nil {
		return delErr.ExtendError(
			"problem deleting file metadata for file " + fileMetadata.Name,
		)
	}

	return nil
}

func (ctrl *Controller) AbortFileMultipartUpload(ctx *gin.Context) {
	apiErr := ctrl.abortFileMultipartUploadProcess(ctx)
	if apiErr != nil {
		_ = ctx.Error(fmt.Errorf("problem processing request: %w", apiErr))

		ctx.JSON(apiErr.statusCode, AbortFileMultipartUploadResponse{apiErr.PublicResponse()})

		return
	}

	ctx.JSON(http.StatusOK, AbortFileMultipartUploadResponse{nil})
}
