package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MultipartFragment struct {
	// Entity tag returned when the part was uploaded.
	ETag string `json:"etag"`
	// Date and time at which the part was uploaded.
	LastModified time.Time `json:"lastModified"`
	// Part number identifying the part. This is a positive integer between 1 and
	// 10,000.
	PartNumber int32 `json:"partNumber"`
	// Size in bytes of the uploaded part data.
	Size int64 `json:"size"`
}

type GetFileMultipartUploadInfoRequest struct {
	FileID   string
	FileETag string
}

type GetFileMultipartUploadInfoResponse struct {
	*FileMetadata
	Parts []MultipartFragment `json:"parts"`
}

func (ctrl *Controller) getFileMultipartUploadInfoParse(ctx *gin.Context) GetFileMultipartUploadInfoRequest {
	etag := ctx.Query("etag")
	if !strings.HasPrefix(etag, "\"") {
		etag = "\"" + etag
	}
	if !strings.HasSuffix(etag, "\"") {
		etag = etag + "\""
	}
	return GetFileMultipartUploadInfoRequest{
		FileID:   ctx.Param("id"),
		FileETag: etag,
	}
}

func (ctrl *Controller) getFileMetadataByETag(
	ctx context.Context,
	fileETag string,
	headers http.Header,
) (FileMetadata, *APIError) {
	fileMetadatas, apiErr := ctrl.metadataStorage.GetFilesByETag(ctx, fileETag, headers)
	if apiErr != nil {
		return FileMetadata{}, apiErr
	}

	if fileMetadatas == nil || len(fileMetadatas) <= 0 {
		errMsg := fmt.Sprintf("can not find file by etag %s", fileETag)
		return FileMetadata{}, BadDataError(errors.New(errMsg), errMsg)
	}

	return fileMetadatas[0], nil
}

func (ctrl *Controller) getFileMultipartUploadInfoProcess(
	ctx *gin.Context,
) (GetFileMultipartUploadInfoResponse, *APIError) {
	req := ctrl.getFileMultipartUploadInfoParse(ctx)

	if req.FileID == "" {
		return GetFileMultipartUploadInfoResponse{}, BadDataError(errors.New("parameter id must be not empty"), "parameter id must be not empty")
	}

	fileMetadata := FileMetadata{}

	if req.FileID == "etag" {
		if req.FileETag == "" {
			return GetFileMultipartUploadInfoResponse{}, BadDataError(errors.New("query parameter etag must be not empty"), "query parameter etag must be not empty")
		}
		fm, apiErr := ctrl.getFileMetadataByETag(
			ctx.Request.Context(), req.FileETag, ctx.Request.Header,
		)
		if apiErr != nil {
			return GetFileMultipartUploadInfoResponse{}, apiErr
		}

		fileMetadata = fm
	} else {
		fm, _, apiErr := ctrl.getFileMetadata(
			ctx.Request.Context(), req.FileID, false, ctx.Request.Header,
		)
		if apiErr != nil {
			return GetFileMultipartUploadInfoResponse{}, apiErr
		}

		fileMetadata = fm
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	resp := GetFileMultipartUploadInfoResponse{&fileMetadata, nil}

	if fileMetadata.UploadID != "" {
		parts, err := ctrl.contentStorage.ListParts(ctx, objectKey, fileMetadata.UploadID)
		if err != nil {
			return GetFileMultipartUploadInfoResponse{}, err
		}
		resp.Parts = parts
	}
	return resp, nil
}

func (ctrl *Controller) GetFileMultipartUploadInfo(ctx *gin.Context) {
	resp, apiErr := ctrl.getFileMultipartUploadInfoProcess(ctx)
	if apiErr != nil {
		_ = ctx.Error(apiErr)

		ctx.JSON(apiErr.statusCode, CommonResponse{
			Code:    apiErr.statusCode,
			Message: apiErr.PublicResponse().Message,
		})

		return
	}

	ctx.JSON(http.StatusOK, CommonResponse{
		http.StatusOK,
		"ok",
		resp,
	})
}
