package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) deleteOrphans(ctx *gin.Context) ([]string, *APIError) {
	toDelete, apiErr := ctrl.listOrphans(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	for _, f := range toDelete {
		if apiErr := ctrl.contentStorage.DeleteFile(ctx, f); apiErr != nil {
			return nil, apiErr
		}
	}

	return toDelete, nil
}

func (ctrl *Controller) DeleteOrphans(ctx *gin.Context) {
	files, apiErr := ctrl.deleteOrphans(ctx)
	if apiErr != nil {
		_ = ctx.Error(fmt.Errorf("problem processing request: %w", apiErr))

		ctx.JSON(apiErr.statusCode, CommonResponse{
			Code:    apiErr.statusCode,
			Message: apiErr.PublicResponse().Message,
		})

		return
	}

	ctx.JSON(
		http.StatusOK,
		CommonResponse{
			http.StatusOK,
			"ok",
			ListOrphansResponse{
				files,
			},
		},
	)
}
