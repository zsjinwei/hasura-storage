package controller

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type UploadFileMultipartWithPresignedURLRequest struct {
	fileID    string
	signature string
	Expires   int
}

func (ctrl *Controller) uploadFileMultipartWithPresignedURLParse(
	ctx *gin.Context,
) (UploadFileMultipartWithPresignedURLRequest, *APIError) {
	expires, apiErr := expiresIn(ctx.Request.URL.Query())
	if apiErr != nil {
		return UploadFileMultipartWithPresignedURLRequest{
			fileID:    "",
			signature: "",
			Expires:   expires,
		}, apiErr //nolint: exhaustruct
	}

	signature := make(url.Values, len(ctx.Request.URL.Query()))
	for k, v := range ctx.Request.URL.Query() {
		switch k {
		default:
			signature[k] = v
		}
	}

	return UploadFileMultipartWithPresignedURLRequest{
		fileID:    ctx.Param("id"),
		signature: signature.Encode(),
		Expires:   expires,
	}, nil
}

func (ctrl *Controller) uploadFileMultipartWithPresignedURL(ctx *gin.Context) *APIError {
	req, apiErr := ctrl.uploadFileMultipartWithPresignedURLParse(ctx)
	if apiErr != nil {
		return apiErr
	}

	fileMetadata, _, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(),
		req.fileID,
		false,
		http.Header{"x-hasura-admin-secret": []string{ctrl.hasuraAdminSecret}},
	)
	if apiErr != nil {
		return apiErr
	}

	objectKey := fileMetadata.ObjectKey
	if objectKey == "" {
		objectKey = fileMetadata.ID
	}

	proxy, apiErr := ctrl.contentStorage.PutFileWithPresignedURL(ctx, objectKey, req.signature,
		ctx.Request.Header)
	if apiErr != nil {
		return apiErr
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)

	return nil
}

func (ctrl *Controller) UploadFileMultipartWithPresignedURL(ctx *gin.Context) {
	apiErr := ctrl.uploadFileMultipartWithPresignedURL(ctx)
	if apiErr != nil {
		_ = ctx.Error(apiErr)

		ctx.JSON(apiErr.statusCode, CommonResponse{
			Code:    apiErr.statusCode,
			Message: apiErr.PublicResponse().Message,
		})

		return
	}
}
