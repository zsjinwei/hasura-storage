package controller

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type downloadFileRequest struct {
	fileID   string
	fileName string
	token    string
	headers  getFileInformationHeaders
}

// Only used if the request fails.
type DownloadFileResponse struct {
	Error *ErrorResponse `json:"error"`
}

func (ctrl *Controller) downloadFileParse(ctx *gin.Context) (downloadFileRequest, *APIError) {
	var headers getFileInformationHeaders
	if err := ctx.ShouldBindHeader(&headers); err != nil {
		return downloadFileRequest{}, InternalServerError(
			fmt.Errorf("problem parsing request headers: %w", err),
		)
	}

	return downloadFileRequest{ctx.Param("id"), ctx.Param("name"), ctx.Query("token"), headers}, nil
}

func (ctrl *Controller) downloadFileProcess(ctx *gin.Context) (*FileResponse, *APIError) {
	req, apiErr := ctrl.downloadFileParse(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	if req.fileName == "" {
		return nil, BadDataError(errors.New("parameter name must be not empty"), "parameter name must be not empty")
	}

	if req.token == "" {
		return nil, BadDataError(errors.New("query parameter token must be not empty"), "query parameter token must be not empty")
	}

	headers := http.Header(ctx.Request.Header)
	headers["x-hasura-admin-secret"] = []string{ctrl.hasuraAdminSecret}

	fileMetadata, bucketMetadata, apiErr := ctrl.getFileMetadata(
		ctx.Request.Context(), req.fileID, true,
		headers,
	)
	if apiErr != nil {
		return nil, apiErr
	}

	if fileMetadata.Name != req.fileName && req.fileName != "default" {
		return nil, BadDataError(errors.New("name is not matched with file"), "name is not matched with file")
	}

	if fileMetadata.ETag != req.token && fileMetadata.ETag != fmt.Sprintf("\"%s\"", req.token) {
		return nil, BadDataError(errors.New("token is invalid"), "token is invalid")
	}

	filePath := fileMetadata.ID
	if len(fileMetadata.ObjectKey) > 0 {
		filePath = fileMetadata.ObjectKey
	}

	downloadFunc := func() (*File, *APIError) {
		return ctrl.contentStorage.GetFile(ctx, filePath, ctx.Request.Header)
	}

	response, apiErr := ctrl.processFileToDownload(
		ctx, downloadFunc, fileMetadata, bucketMetadata.CacheControl, &req.headers)
	if apiErr != nil {
		return nil, apiErr
	}

	if response.statusCode == http.StatusOK {
		// if we want to download files at some point prepend `attachment;` before filename
		response.headers.Add(
			"Content-Disposition",
			fmt.Sprintf(`inline; filename="%s"`, url.QueryEscape(fileMetadata.Name)),
		)
	}

	return response, nil
}

func (ctrl *Controller) DownloadFile(ctx *gin.Context) {
	response, apiErr := ctrl.downloadFileProcess(ctx)
	if apiErr != nil {
		_ = ctx.Error(apiErr)

		ctx.JSON(apiErr.statusCode, DownloadFileResponse{apiErr.PublicResponse()})

		return
	}

	defer response.body.Close()

	response.Write(ctx)
}
