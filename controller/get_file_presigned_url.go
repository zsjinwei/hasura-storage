package controller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetFilePresignedURLResponse struct {
	URL        string `json:"url,omitempty"`
	Expiration int    `json:"expiration,omitempty"`
}

type GetFilePresignedURLRequest struct {
	FileID string
}

func (ctrl *Controller) getFilePresignedURLParse(ctx *gin.Context) GetFilePresignedURLRequest {
	return GetFilePresignedURLRequest{
		FileID: ctx.Param("id"),
	}
}

func (ctrl *Controller) getFilePresignedURL(
	ctx *gin.Context,
) (GetFilePresignedURLResponse, *APIError) {
	req := ctrl.getFilePresignedURLParse(ctx)

	fileMetadata, bucketMetadata, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(), req.FileID, true, ctx.Request.Header,
	)
	if apiErr != nil {
		return GetFilePresignedURLResponse{}, apiErr
	}

	if !bucketMetadata.PresignedURLsEnabled {
		err := errors.New( //nolint: goerr113
			"presigned URLs are not enabled on the bucket where this file is located in",
		)
		return GetFilePresignedURLResponse{}, ForbiddenError(err, err.Error())
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	signature, apiErr := ctrl.contentStorage.CreateGetObjectPresignedURL(
		ctx,
		objectKey,
		time.Duration(bucketMetadata.DownloadExpiration)*time.Second,
	)
	if apiErr != nil {
		return GetFilePresignedURLResponse{},
			apiErr.ExtendError(
				"problem creating presigned URL for file " + fileMetadata.Name,
			)
	}

	url := fmt.Sprintf(
		"%s%s/files/%s/presignedurl/content?%s",
		ctrl.publicURL, ctrl.apiRootPrefix, fileMetadata.ID, signature,
	)
	return GetFilePresignedURLResponse{url, bucketMetadata.DownloadExpiration}, nil
}

func (ctrl *Controller) GetFilePresignedURL(ctx *gin.Context) {
	resp, apiErr := ctrl.getFilePresignedURL(ctx)
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
