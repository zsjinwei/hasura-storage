package metadata

import (
	"context"
	"errors"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/nhost/hasura-storage/controller"
)

func ptr[T any](x T) *T {
	return &x
}

func parseGraphqlError(err error) *controller.APIError {
	var ghErr *clientv2.ErrorResponse
	if errors.As(err, &ghErr) {
		code, ok := (*ghErr.GqlErrors)[0].Extensions["code"]
		if !ok {
			return controller.InternalServerError(err)
		}
		switch code {
		case "access-denied", "validation-failed", "permission-error":
			return controller.ForbiddenError(ghErr, "you are not authorized")
		case "data-exception", "constraint-violation":
			return controller.BadDataError(err, ghErr.Error())
		default:
			return controller.InternalServerError(err)
		}
	}

	return controller.InternalServerError(err)
}

func (md *FileMetadataSummaryFragment) ToControllerType() controller.FileSummary {
	return controller.FileSummary{
		ID:         md.GetID(),
		Name:       *md.GetName(),
		BucketID:   md.GetBucketID(),
		IsUploaded: *md.GetIsUploaded(),
	}
}

func (md *BucketMetadataFragment) ToControllerType() controller.BucketMetadata {
	return controller.BucketMetadata{
		ID:                   md.GetID(),
		MinUploadFile:        int(md.GetMinUploadFileSize()),
		MaxUploadFile:        int(md.GetMaxUploadFileSize()),
		PresignedURLsEnabled: md.GetPresignedUrlsEnabled(),
		DownloadExpiration:   int(md.GetDownloadExpiration()),
		CreatedAt:            md.GetCreatedAt(),
		UpdatedAt:            md.GetUpdatedAt(),
		CacheControl:         *md.GetCacheControl(),
		UploadExpiration:     int(md.GetUploadExpiration()),
	}
}

func (md *FileMetadataFragment) ToControllerType() controller.FileMetadata {
	ID := md.GetID()

	Name := ""
	if md.GetName() != nil {
		Name = *md.GetName()
	}

	Size := int64(0)
	if md.GetSize() != nil {
		Size = *md.GetSize()
	}

	BucketID := md.GetBucketID()

	ETag := ""
	if md.GetEtag() != nil {
		ETag = *md.GetEtag()
	}

	CreatedAt := md.GetCreatedAt()

	UpdatedAt := md.GetUpdatedAt()

	IsUploaded := false
	if md.GetIsUploaded() != nil {
		IsUploaded = *md.GetIsUploaded()
	}

	MimeType := ""
	if md.GetMimeType() != nil {
		MimeType = *md.GetMimeType()
	}

	Metadata := md.GetMetadata()

	ObjectKey := ""
	if md.GetObjectKey() != nil {
		ObjectKey = *md.GetObjectKey()
	}

	ChunkSize := int64(0)
	if md.GetChunkSize() != nil {
		ChunkSize = *md.GetChunkSize()
	}

	ChunkCount := int64(0)
	if md.GetChunkCount() != nil {
		ChunkCount = *md.GetChunkCount()
	}

	UploadID := ""
	if md.GetUploadID() != nil {
		UploadID = *md.GetUploadID()
	}

	return controller.FileMetadata{
		ID:         ID,
		Name:       Name,
		Size:       Size,
		BucketID:   BucketID,
		ETag:       ETag,
		CreatedAt:  CreatedAt,
		UpdatedAt:  UpdatedAt,
		IsUploaded: IsUploaded,
		MimeType:   MimeType,
		Metadata:   Metadata,
		ObjectKey:  ObjectKey,
		ChunkSize:  ChunkSize,
		ChunkCount: ChunkCount,
		UploadID:   UploadID,
	}
}

func WithHeaders(header http.Header) clientv2.RequestInterceptor {
	return func(
		ctx context.Context,
		req *http.Request,
		gqlInfo *clientv2.GQLRequestInfo,
		res interface{},
		next clientv2.RequestInterceptorFunc,
	) error {
		for k, v := range header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
		return next(ctx, req, gqlInfo, res)
	}
}

type Hasura struct {
	cl *Client
}

func NewHasura(endpoint string) *Hasura {
	return &Hasura{
		cl: NewClient(
			&http.Client{},
			endpoint,
			&clientv2.Options{},
		),
	}
}

func (h *Hasura) GetBucketByID(
	ctx context.Context,
	bucketID string,
	headers http.Header,
) (controller.BucketMetadata, *controller.APIError) {
	resp, err := h.cl.GetBucket(
		ctx,
		bucketID,
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return controller.BucketMetadata{}, aerr.ExtendError("problem getting bucket metadata")
	}

	if resp.Bucket == nil || resp.Bucket.ID == "" {
		return controller.BucketMetadata{}, controller.ErrBucketNotFound
	}

	return resp.Bucket.ToControllerType(), nil
}

func (h *Hasura) InitializeFile(
	ctx context.Context,
	fileID, name string, size int64, bucketID, mimeType string,
	objectKey string, chunkSize int64, chunkCount int64, uploadId string,
	headers http.Header,
) *controller.APIError {
	if objectKey == "" {
		objectKey = fileID
	}
	if chunkSize <= 0 || chunkCount <= 0 {
		chunkSize = size
		chunkCount = 1
	}
	_, err := h.cl.InsertFile(
		ctx,
		FilesInsertInput{
			BucketID:   ptr(bucketID),
			ID:         ptr(fileID),
			MimeType:   ptr(mimeType),
			Name:       ptr(name),
			Size:       ptr(size),
			ObjectKey:  ptr(objectKey),
			ChunkSize:  ptr(chunkSize),
			ChunkCount: ptr(chunkCount),
			UploadID:   ptr(uploadId),
		},
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return aerr.ExtendError("problem initializing file metadata")
	}

	return nil
}

func (h *Hasura) PopulateMetadata(
	ctx context.Context,
	fileID, name string, size int64, bucketID, etag string, isUploaded bool, mimeType string,
	objectKey string, chunkSize int64, chunkCount int64, uploadId string,
	metadata map[string]any,
	headers http.Header,
) (controller.FileMetadata, *controller.APIError) {
	if objectKey == "" {
		objectKey = fileID
	}
	if chunkSize <= 0 || chunkCount <= 0 {
		chunkSize = size
		chunkCount = 1
	}
	resp, err := h.cl.UpdateFile(
		ctx,
		fileID,
		FilesSetInput{
			BucketID:   ptr(bucketID),
			Etag:       ptr(etag),
			IsUploaded: ptr(isUploaded),
			Metadata:   metadata,
			MimeType:   ptr(mimeType),
			Name:       ptr(name),
			Size:       ptr(size),
			ObjectKey:  ptr(objectKey),
			ChunkSize:  ptr(chunkSize),
			ChunkCount: ptr(chunkCount),
			UploadID:   ptr(uploadId),
		},
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return controller.FileMetadata{}, aerr.ExtendError("problem populating file metadata")
	}

	if resp.UpdateFile == nil || resp.UpdateFile.ID == "" {
		return controller.FileMetadata{}, controller.ErrFileNotFound
	}

	return resp.UpdateFile.ToControllerType(), nil
}

func (h *Hasura) GetFileByID(
	ctx context.Context,
	fileID string,
	headers http.Header,
) (controller.FileMetadata, *controller.APIError) {
	resp, err := h.cl.GetFile(
		ctx,
		fileID,
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return controller.FileMetadata{}, aerr.ExtendError("problem getting file metadata")
	}

	if resp.File == nil || resp.File.ID == "" {
		return controller.FileMetadata{}, controller.ErrFileNotFound
	}

	return resp.File.ToControllerType(), nil
}

func (h *Hasura) GetFilesByETag(
	ctx context.Context,
	fileETag string,
	headers http.Header,
) ([]controller.FileMetadata, *controller.APIError) {
	resp, err := h.cl.GetFilesByETag(
		ctx,
		fileETag,
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return nil, aerr.ExtendError("problem listing files")
	}

	files := make([]controller.FileMetadata, len(resp.Files))
	for i, f := range resp.Files {
		files[i] = f.ToControllerType()
	}

	return files, nil
}

func (h *Hasura) SetIsUploaded(
	ctx context.Context, fileID string, isUploaded bool, headers http.Header,
) *controller.APIError {
	resp, err := h.cl.UpdateFile(
		ctx,
		fileID,
		FilesSetInput{
			IsUploaded: ptr(isUploaded),
		},
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return aerr.ExtendError("problem setting file as uploaded")
	}

	if resp.UpdateFile == nil || resp.UpdateFile.ID == "" {
		return controller.ErrFileNotFound
	}

	return nil
}

func (h *Hasura) DeleteFileByID(
	ctx context.Context,
	fileID string,
	headers http.Header,
) *controller.APIError {
	resp, err := h.cl.DeleteFile(
		ctx,
		fileID,
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return aerr.ExtendError("problem deleting file")
	}

	if resp.DeleteFile == nil || resp.DeleteFile.ID == "" {
		return controller.ErrFileNotFound
	}

	return nil
}

func (h *Hasura) ListFiles(
	ctx context.Context,
	headers http.Header,
) ([]controller.FileSummary, *controller.APIError) {
	resp, err := h.cl.ListFilesSummary(
		ctx,
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return nil, aerr.ExtendError("problem listing files")
	}

	files := make([]controller.FileSummary, len(resp.Files))
	for i, f := range resp.Files {
		files[i] = f.ToControllerType()
	}

	return files, nil
}

func (h *Hasura) InsertVirus(
	ctx context.Context,
	fileID, filename, virus string,
	userSession map[string]any,
	headers http.Header,
) *controller.APIError {
	_, err := h.cl.InsertVirus(
		ctx,
		VirusInsertInput{
			FileID:      ptr(fileID),
			Filename:    ptr(filename),
			UserSession: userSession,
			Virus:       ptr(virus),
		},
		WithHeaders(headers),
	)
	if err != nil {
		aerr := parseGraphqlError(err)
		return aerr.ExtendError("problem inserting virus")
	}

	return nil
}
