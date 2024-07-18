package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetFileMultipartUploadPresignedURLResponse struct {
	URLs       []string `json:"urls,omitempty"`
	Expiration int      `json:"expiration,omitempty"`
}

type GetFileMultipartUploadPresignedURLRequest struct {
	FileID string
}

func (ctrl *Controller) getFileMultipartUploadPresignedURLParse(ctx *gin.Context) GetFileMultipartUploadPresignedURLRequest {
	return GetFileMultipartUploadPresignedURLRequest{
		FileID: ctx.Param("id"),
	}
}

func (ctrl *Controller) getFileMultipartUploadPresignedURL(
	ctx *gin.Context,
) (GetFileMultipartUploadPresignedURLResponse, *APIError) {
	req := ctrl.getFileMultipartUploadPresignedURLParse(ctx)

	fileMetadata, bucketMetadata, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(), req.FileID, false, ctx.Request.Header,
	)
	if apiErr != nil {
		return GetFileMultipartUploadPresignedURLResponse{}, apiErr
	}

	if !bucketMetadata.PresignedURLsEnabled {
		err := errors.New( //nolint: goerr113
			"presigned URLs are not enabled on the bucket where this file is located in",
		)
		return GetFileMultipartUploadPresignedURLResponse{}, ForbiddenError(err, err.Error())
	}

	if fileMetadata.UploadID == "" {
		return GetFileMultipartUploadPresignedURLResponse{},
			apiErr.ExtendError(
				"no multipart uploadId for file " + fileMetadata.Name,
			)
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	urls := make([]string, 0, fileMetadata.ChunkCount)

	for i := int64(1); i <= fileMetadata.ChunkCount; i++ {
		signature, apiErr := ctrl.contentStorage.CreateUploadPartPresignedURL(
			ctx,
			objectKey,
			fileMetadata.UploadID,
			int32(i),
			time.Duration(bucketMetadata.UploadExpiration)*time.Second,
		)
		if apiErr != nil {
			return GetFileMultipartUploadPresignedURLResponse{},
				apiErr.ExtendError(
					"problem creating presigned URL for file " + fileMetadata.Name,
				)
		}

		url := fmt.Sprintf(
			"%s%s/files/%s/multipart/presignedurl/content?%s",
			ctrl.publicURL, ctrl.apiRootPrefix, fileMetadata.ID, signature,
		)
		urls = append(urls, url)
	}

	return GetFileMultipartUploadPresignedURLResponse{urls, bucketMetadata.UploadExpiration}, nil
}

func (ctrl *Controller) GetFileMultipartPresignedURL(ctx *gin.Context) {
	resp, apiErr := ctrl.getFileMultipartUploadPresignedURL(ctx)
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
