//go:generate mockgen -destination mock/controller.go -package mock -source=controller.go MetadataStorage
package controller

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nhost/hasura-storage/image"
	"github.com/sirupsen/logrus"
)

type FileSummary struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsUploaded bool   `json:"isUploaded"`
	BucketID   string `json:"bucketId"`
}

type BucketMetadata struct {
	ID                   string
	MinUploadFile        int
	MaxUploadFile        int
	PresignedURLsEnabled bool
	DownloadExpiration   int
	CreatedAt            string
	UpdatedAt            string
	CacheControl         string
	UploadExpiration     int
}

type FileMetadata struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Size             int64          `json:"size"`
	BucketID         string         `json:"bucketId"`
	ETag             string         `json:"etag"`
	CreatedAt        string         `json:"createdAt"`
	UpdatedAt        string         `json:"updatedAt"`
	IsUploaded       bool           `json:"isUploaded"`
	MimeType         string         `json:"mimeType"`
	UploadedByUserID string         `json:"uploadedByUserId"`
	Metadata         map[string]any `json:"metadata"`
	ObjectKey        string         `json:"objectKey"`
	ChunkSize        int64          `json:"chunkSize"`
	ChunkCount       int64          `json:"chunkCount"`
	UploadID         string         `json:"uploadId"`
}

type MetadataStorage interface {
	GetBucketByID(ctx context.Context, id string, headers http.Header) (BucketMetadata, *APIError)
	GetFileByID(ctx context.Context, id string, headers http.Header) (FileMetadata, *APIError)
	GetFilesByETag(ctx context.Context, etag string, headers http.Header) ([]FileMetadata, *APIError)
	InitializeFile(
		ctx context.Context,
		id, name string, size int64, bucketID, mimeType string,
		objectKey string, chunkSize int64, chunkCount int64, uploadId string,
		headers http.Header,
	) *APIError
	PopulateMetadata(
		ctx context.Context,
		id, name string, size int64, bucketID, etag string, IsUploaded bool, mimeType string,
		objectKey string, chunkSize int64, chunkCount int64, uploadId string,
		metadata map[string]any,
		headers http.Header) (FileMetadata, *APIError,
	)
	SetIsUploaded(
		ctx context.Context,
		fileID string,
		isUploaded bool,
		headers http.Header,
	) *APIError
	DeleteFileByID(ctx context.Context, fileID string, headers http.Header) *APIError
	ListFiles(ctx context.Context, headers http.Header) ([]FileSummary, *APIError)
	InsertVirus(
		ctx context.Context,
		fileID, filename, virus string,
		userSession map[string]any,
		headers http.Header,
	) *APIError
}

type ContentStorage interface {
	PutFile(
		ctx context.Context,
		content io.ReadSeeker,
		filepath, contentType string,
	) (string, *APIError)
	GetFile(ctx context.Context, filepath string, headers http.Header) (*File, *APIError)
	CreateGetObjectPresignedURL(
		ctx context.Context,
		filepath string,
		expire time.Duration,
	) (string, *APIError)
	GetFileWithPresignedURL(
		ctx context.Context, filepath, signature string, headers http.Header,
	) (*File, *APIError)
	PutFileWithPresignedURL(
		ctx context.Context, filepath, signature string, headers http.Header,
	) (*httputil.ReverseProxy, *APIError)
	DeleteFile(ctx context.Context, filepath string) *APIError
	ListFiles(ctx context.Context) ([]string, *APIError)
	CreateMultipartUpload(
		ctx context.Context,
		filepath string,
		contentType string,
	) (string, *APIError)
	ListParts(
		ctx context.Context,
		filepath string,
		uploadId string,
	) ([]MultipartFragment, *APIError)
	UploadPart(
		ctx context.Context,
		filepath string,
		uploadId string,
		partNumber int32,
		body io.ReadSeeker,
	) (string, *APIError)
	CreatePutObjectPresignedURL(
		ctx context.Context,
		filepath string,
		contentType string,
		expire time.Duration,
	) (string, *APIError)
	CreateUploadPartPresignedURL(
		ctx context.Context,
		filepath string,
		uploadId string,
		partNumber int32,
		expire time.Duration,
	) (string, *APIError)
	CompleteMultipartUpload(
		ctx context.Context,
		filepath string,
		uploadId string,
	) (string, *APIError)
	AbortMultipartUpload(
		ctx context.Context,
		filepath string,
		uploadId string,
	) *APIError
}

type Antivirus interface {
	ScanReader(r io.ReaderAt) *APIError
}

type Controller struct {
	publicURL         string
	apiRootPrefix     string
	hasuraAdminSecret string
	metadataStorage   MetadataStorage
	contentStorage    ContentStorage
	imageTransformer  *image.Transformer
	av                Antivirus
	logger            *logrus.Logger
}

func New(
	publicURL string,
	apiRootPrefix string,
	hasuraAdminSecret string,
	metadataStorage MetadataStorage,
	contentStorage ContentStorage,
	imageTransformer *image.Transformer,
	av Antivirus,
	logger *logrus.Logger,
) *Controller {
	return &Controller{
		publicURL,
		apiRootPrefix,
		hasuraAdminSecret,
		metadataStorage,
		contentStorage,
		imageTransformer,
		av,
		logger,
	}
}

func corsConfig(allowedOrigins []string) cors.Config {
	return cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "PUT", "POST", "HEAD", "DELETE"},
		AllowHeaders: []string{
			"Authorization", "Origin", "if-match", "if-none-match", "if-modified-since", "if-unmodified-since",
			"x-hasura-admin-secret", "x-nhost-bucket-id", "x-nhost-file-name", "x-nhost-file-id",
			"x-hasura-role",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type", "Cache-Control", "ETag", "Last-Modified", "X-Error",
		},
		MaxAge: 12 * time.Hour, //nolint: mnd
	}
}

func (ctrl *Controller) SetupRouter(
	trustedProxies []string,
	apiRootPrefix string,
	corsOrigins []string,
	corsAllowCredentials bool,
	middleware ...gin.HandlerFunc,
) (*gin.Engine, error) {
	router := gin.New()
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		return nil, fmt.Errorf("problem setting trusted proxies: %w", err)
	}

	// lower values make uploads slower but keeps service memory usage low
	router.MaxMultipartMemory = 1 << 20 //nolint:mnd  // 1 MB
	router.Use(gin.Recovery())

	for _, mw := range middleware {
		router.Use(mw)
	}

	corsConfig := corsConfig(corsOrigins)
	if corsAllowCredentials {
		corsConfig.AllowCredentials = true
	}

	router.Use(cors.New(corsConfig))

	router.GET("/healthz", ctrl.Health)

	apiRoot := router.Group(apiRootPrefix)
	{
		apiRoot.GET("/openapi.yaml", ctrl.OpenAPI)
		apiRoot.GET("/version", ctrl.Version)
	}
	files := apiRoot.Group("/files")
	{
		files.POST("", ctrl.UploadFile) // To delete
		files.POST("/", ctrl.UploadFile)
		files.GET("/:id", ctrl.GetFile)
		files.HEAD("/:id", ctrl.GetFileInformation)
		files.PUT("/:id", ctrl.UpdateFile)
		files.DELETE("/:id", ctrl.DeleteFile)
		files.GET("/:id/presignedurl", ctrl.GetFilePresignedURL)
		files.GET("/:id/presignedurl/content", ctrl.GetFileWithPresignedURL)
		files.GET("/:id/download/:name", ctrl.DownloadFile)
		files.GET("/:id/multipart", ctrl.GetFileMultipartUploadInfo)
		files.GET("/:id/multipart/presignedurl", ctrl.GetFileMultipartPresignedURL)
		files.PUT("/:id/multipart/presignedurl/content", ctrl.UploadFileMultipartWithPresignedURL)
		files.POST("/multipart", ctrl.CreateFileMultipartUpload)
		files.POST("/:id/multipart/complete", ctrl.CompleteFileMultipartUpload)
		files.POST("/:id/multipart/abort", ctrl.AbortFileMultipartUpload)
	}

	ops := apiRoot.Group("/ops")
	{
		ops.POST("list-orphans", ctrl.ListOrphans)
		ops.POST("delete-orphans", ctrl.DeleteOrphans)
		ops.POST("list-broken-metadata", ctrl.ListBrokenMetadata)
		ops.POST("delete-broken-metadata", ctrl.DeleteBrokenMetadata)
		ops.POST("list-not-uploaded", ctrl.ListNotUploaded)
	}
	return router, nil
}

func (ctrl *Controller) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"healthz": "ok",
	})
}
