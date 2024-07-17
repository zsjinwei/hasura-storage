// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package metadata

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
)

type Client struct {
	Client *clientv2.Client
}

func NewClient(cli *http.Client, baseURL string, options *clientv2.Options, interceptors ...clientv2.RequestInterceptor) *Client {
	return &Client{Client: clientv2.NewClient(cli, baseURL, options, interceptors...)}
}

type QueryRoot struct {
	Bucket           *Buckets         "json:\"bucket,omitempty\" graphql:\"bucket\""
	Buckets          []*Buckets       "json:\"buckets\" graphql:\"buckets\""
	BucketsAggregate BucketsAggregate "json:\"bucketsAggregate\" graphql:\"bucketsAggregate\""
	File             *Files           "json:\"file,omitempty\" graphql:\"file\""
	Files            []*Files         "json:\"files\" graphql:\"files\""
	FilesAggregate   FilesAggregate   "json:\"filesAggregate\" graphql:\"filesAggregate\""
	Virus            *Virus           "json:\"virus,omitempty\" graphql:\"virus\""
	Viruses          []*Virus         "json:\"viruses\" graphql:\"viruses\""
	VirusesAggregate VirusAggregate   "json:\"virusesAggregate\" graphql:\"virusesAggregate\""
}
type MutationRoot struct {
	DeleteBucket      *Buckets                   "json:\"deleteBucket,omitempty\" graphql:\"deleteBucket\""
	DeleteBuckets     *BucketsMutationResponse   "json:\"deleteBuckets,omitempty\" graphql:\"deleteBuckets\""
	DeleteFile        *Files                     "json:\"deleteFile,omitempty\" graphql:\"deleteFile\""
	DeleteFiles       *FilesMutationResponse     "json:\"deleteFiles,omitempty\" graphql:\"deleteFiles\""
	DeleteVirus       *Virus                     "json:\"deleteVirus,omitempty\" graphql:\"deleteVirus\""
	DeleteViruses     *VirusMutationResponse     "json:\"deleteViruses,omitempty\" graphql:\"deleteViruses\""
	InsertBucket      *Buckets                   "json:\"insertBucket,omitempty\" graphql:\"insertBucket\""
	InsertBuckets     *BucketsMutationResponse   "json:\"insertBuckets,omitempty\" graphql:\"insertBuckets\""
	InsertFile        *Files                     "json:\"insertFile,omitempty\" graphql:\"insertFile\""
	InsertFiles       *FilesMutationResponse     "json:\"insertFiles,omitempty\" graphql:\"insertFiles\""
	InsertVirus       *Virus                     "json:\"insertVirus,omitempty\" graphql:\"insertVirus\""
	InsertViruses     *VirusMutationResponse     "json:\"insertViruses,omitempty\" graphql:\"insertViruses\""
	UpdateBucket      *Buckets                   "json:\"updateBucket,omitempty\" graphql:\"updateBucket\""
	UpdateBuckets     *BucketsMutationResponse   "json:\"updateBuckets,omitempty\" graphql:\"updateBuckets\""
	UpdateFile        *Files                     "json:\"updateFile,omitempty\" graphql:\"updateFile\""
	UpdateFiles       *FilesMutationResponse     "json:\"updateFiles,omitempty\" graphql:\"updateFiles\""
	UpdateVirus       *Virus                     "json:\"updateVirus,omitempty\" graphql:\"updateVirus\""
	UpdateViruses     *VirusMutationResponse     "json:\"updateViruses,omitempty\" graphql:\"updateViruses\""
	UpdateBucketsMany []*BucketsMutationResponse "json:\"update_buckets_many,omitempty\" graphql:\"update_buckets_many\""
	UpdateFilesMany   []*FilesMutationResponse   "json:\"update_files_many,omitempty\" graphql:\"update_files_many\""
	UpdateVirusMany   []*VirusMutationResponse   "json:\"update_virus_many,omitempty\" graphql:\"update_virus_many\""
}
type FileMetadataFragment struct {
	ID               string                 "json:\"id\" graphql:\"id\""
	Name             *string                "json:\"name,omitempty\" graphql:\"name\""
	Size             *int64                 "json:\"size,omitempty\" graphql:\"size\""
	BucketID         string                 "json:\"bucketId\" graphql:\"bucketId\""
	Etag             *string                "json:\"etag,omitempty\" graphql:\"etag\""
	CreatedAt        string                 "json:\"createdAt\" graphql:\"createdAt\""
	UpdatedAt        string                 "json:\"updatedAt\" graphql:\"updatedAt\""
	IsUploaded       *bool                  "json:\"isUploaded,omitempty\" graphql:\"isUploaded\""
	MimeType         *string                "json:\"mimeType,omitempty\" graphql:\"mimeType\""
	UploadedByUserID *string                "json:\"uploadedByUserId,omitempty\" graphql:\"uploadedByUserId\""
	Metadata         map[string]interface{} "json:\"metadata,omitempty\" graphql:\"metadata\""
	ObjectKey        *string                "json:\"objectKey,omitempty\" graphql:\"objectKey\""
	ChunkSize        *int64                 "json:\"chunkSize,omitempty\" graphql:\"chunkSize\""
	ChunkCount       *int64                 "json:\"chunkCount,omitempty\" graphql:\"chunkCount\""
	UploadID         *string                "json:\"uploadId,omitempty\" graphql:\"uploadId\""
}

func (t *FileMetadataFragment) GetID() string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.ID
}
func (t *FileMetadataFragment) GetName() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.Name
}
func (t *FileMetadataFragment) GetSize() *int64 {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.Size
}
func (t *FileMetadataFragment) GetBucketID() string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.BucketID
}
func (t *FileMetadataFragment) GetEtag() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.Etag
}
func (t *FileMetadataFragment) GetCreatedAt() string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.CreatedAt
}
func (t *FileMetadataFragment) GetUpdatedAt() string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.UpdatedAt
}
func (t *FileMetadataFragment) GetIsUploaded() *bool {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.IsUploaded
}
func (t *FileMetadataFragment) GetMimeType() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.MimeType
}
func (t *FileMetadataFragment) GetUploadedByUserID() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.UploadedByUserID
}
func (t *FileMetadataFragment) GetMetadata() map[string]interface{} {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.Metadata
}
func (t *FileMetadataFragment) GetObjectKey() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.ObjectKey
}
func (t *FileMetadataFragment) GetChunkSize() *int64 {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.ChunkSize
}
func (t *FileMetadataFragment) GetChunkCount() *int64 {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.ChunkCount
}
func (t *FileMetadataFragment) GetUploadID() *string {
	if t == nil {
		t = &FileMetadataFragment{}
	}
	return t.UploadID
}

type FileMetadataSummaryFragment struct {
	ID         string  "json:\"id\" graphql:\"id\""
	Name       *string "json:\"name,omitempty\" graphql:\"name\""
	BucketID   string  "json:\"bucketId\" graphql:\"bucketId\""
	IsUploaded *bool   "json:\"isUploaded,omitempty\" graphql:\"isUploaded\""
}

func (t *FileMetadataSummaryFragment) GetID() string {
	if t == nil {
		t = &FileMetadataSummaryFragment{}
	}
	return t.ID
}
func (t *FileMetadataSummaryFragment) GetName() *string {
	if t == nil {
		t = &FileMetadataSummaryFragment{}
	}
	return t.Name
}
func (t *FileMetadataSummaryFragment) GetBucketID() string {
	if t == nil {
		t = &FileMetadataSummaryFragment{}
	}
	return t.BucketID
}
func (t *FileMetadataSummaryFragment) GetIsUploaded() *bool {
	if t == nil {
		t = &FileMetadataSummaryFragment{}
	}
	return t.IsUploaded
}

type BucketMetadataFragment struct {
	ID                   string  "json:\"id\" graphql:\"id\""
	MinUploadFileSize    int64   "json:\"minUploadFileSize\" graphql:\"minUploadFileSize\""
	MaxUploadFileSize    int64   "json:\"maxUploadFileSize\" graphql:\"maxUploadFileSize\""
	PresignedUrlsEnabled bool    "json:\"presignedUrlsEnabled\" graphql:\"presignedUrlsEnabled\""
	DownloadExpiration   int64   "json:\"downloadExpiration\" graphql:\"downloadExpiration\""
	CreatedAt            string  "json:\"createdAt\" graphql:\"createdAt\""
	UpdatedAt            string  "json:\"updatedAt\" graphql:\"updatedAt\""
	CacheControl         *string "json:\"cacheControl,omitempty\" graphql:\"cacheControl\""
	UploadExpiration     int64   "json:\"uploadExpiration\" graphql:\"uploadExpiration\""
}

func (t *BucketMetadataFragment) GetID() string {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.ID
}
func (t *BucketMetadataFragment) GetMinUploadFileSize() int64 {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.MinUploadFileSize
}
func (t *BucketMetadataFragment) GetMaxUploadFileSize() int64 {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.MaxUploadFileSize
}
func (t *BucketMetadataFragment) GetPresignedUrlsEnabled() bool {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.PresignedUrlsEnabled
}
func (t *BucketMetadataFragment) GetDownloadExpiration() int64 {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.DownloadExpiration
}
func (t *BucketMetadataFragment) GetCreatedAt() string {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.CreatedAt
}
func (t *BucketMetadataFragment) GetUpdatedAt() string {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.UpdatedAt
}
func (t *BucketMetadataFragment) GetCacheControl() *string {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.CacheControl
}
func (t *BucketMetadataFragment) GetUploadExpiration() int64 {
	if t == nil {
		t = &BucketMetadataFragment{}
	}
	return t.UploadExpiration
}

type InsertFile_InsertFile struct {
	ID string "json:\"id\" graphql:\"id\""
}

func (t *InsertFile_InsertFile) GetID() string {
	if t == nil {
		t = &InsertFile_InsertFile{}
	}
	return t.ID
}

type DeleteFile_DeleteFile struct {
	ID string "json:\"id\" graphql:\"id\""
}

func (t *DeleteFile_DeleteFile) GetID() string {
	if t == nil {
		t = &DeleteFile_DeleteFile{}
	}
	return t.ID
}

type InsertVirus_InsertVirus struct {
	ID string "json:\"id\" graphql:\"id\""
}

func (t *InsertVirus_InsertVirus) GetID() string {
	if t == nil {
		t = &InsertVirus_InsertVirus{}
	}
	return t.ID
}

type GetBucket struct {
	Bucket *BucketMetadataFragment "json:\"bucket,omitempty\" graphql:\"bucket\""
}

func (t *GetBucket) GetBucket() *BucketMetadataFragment {
	if t == nil {
		t = &GetBucket{}
	}
	return t.Bucket
}

type GetFile struct {
	File *FileMetadataFragment "json:\"file,omitempty\" graphql:\"file\""
}

func (t *GetFile) GetFile() *FileMetadataFragment {
	if t == nil {
		t = &GetFile{}
	}
	return t.File
}

type GetFilesByETag struct {
	Files []*FileMetadataFragment "json:\"files\" graphql:\"files\""
}

func (t *GetFilesByETag) GetFiles() []*FileMetadataFragment {
	if t == nil {
		t = &GetFilesByETag{}
	}
	return t.Files
}

type ListFilesSummary struct {
	Files []*FileMetadataSummaryFragment "json:\"files\" graphql:\"files\""
}

func (t *ListFilesSummary) GetFiles() []*FileMetadataSummaryFragment {
	if t == nil {
		t = &ListFilesSummary{}
	}
	return t.Files
}

type InsertFile struct {
	InsertFile *InsertFile_InsertFile "json:\"insertFile,omitempty\" graphql:\"insertFile\""
}

func (t *InsertFile) GetInsertFile() *InsertFile_InsertFile {
	if t == nil {
		t = &InsertFile{}
	}
	return t.InsertFile
}

type UpdateFile struct {
	UpdateFile *FileMetadataFragment "json:\"updateFile,omitempty\" graphql:\"updateFile\""
}

func (t *UpdateFile) GetUpdateFile() *FileMetadataFragment {
	if t == nil {
		t = &UpdateFile{}
	}
	return t.UpdateFile
}

type DeleteFile struct {
	DeleteFile *DeleteFile_DeleteFile "json:\"deleteFile,omitempty\" graphql:\"deleteFile\""
}

func (t *DeleteFile) GetDeleteFile() *DeleteFile_DeleteFile {
	if t == nil {
		t = &DeleteFile{}
	}
	return t.DeleteFile
}

type InsertVirus struct {
	InsertVirus *InsertVirus_InsertVirus "json:\"insertVirus,omitempty\" graphql:\"insertVirus\""
}

func (t *InsertVirus) GetInsertVirus() *InsertVirus_InsertVirus {
	if t == nil {
		t = &InsertVirus{}
	}
	return t.InsertVirus
}

const GetBucketDocument = `query GetBucket ($id: String!) {
	bucket(id: $id) {
		... BucketMetadataFragment
	}
}
fragment BucketMetadataFragment on buckets {
	id
	minUploadFileSize
	maxUploadFileSize
	presignedUrlsEnabled
	downloadExpiration
	createdAt
	updatedAt
	cacheControl
	uploadExpiration
}
`

func (c *Client) GetBucket(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*GetBucket, error) {
	vars := map[string]any{
		"id": id,
	}

	var res GetBucket
	if err := c.Client.Post(ctx, "GetBucket", GetBucketDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const GetFileDocument = `query GetFile ($id: uuid!) {
	file(id: $id) {
		... FileMetadataFragment
	}
}
fragment FileMetadataFragment on files {
	id
	name
	size
	bucketId
	etag
	createdAt
	updatedAt
	isUploaded
	mimeType
	uploadedByUserId
	metadata
	objectKey
	chunkSize
	chunkCount
	uploadId
}
`

func (c *Client) GetFile(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*GetFile, error) {
	vars := map[string]any{
		"id": id,
	}

	var res GetFile
	if err := c.Client.Post(ctx, "GetFile", GetFileDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const GetFilesByETagDocument = `query GetFilesByETag ($etag: String!) {
	files(where: {etag:{_eq:$etag}}, order_by: {createdAt:desc}) {
		... FileMetadataFragment
	}
}
fragment FileMetadataFragment on files {
	id
	name
	size
	bucketId
	etag
	createdAt
	updatedAt
	isUploaded
	mimeType
	uploadedByUserId
	metadata
	objectKey
	chunkSize
	chunkCount
	uploadId
}
`

func (c *Client) GetFilesByETag(ctx context.Context, etag string, interceptors ...clientv2.RequestInterceptor) (*GetFilesByETag, error) {
	vars := map[string]any{
		"etag": etag,
	}

	var res GetFilesByETag
	if err := c.Client.Post(ctx, "GetFilesByETag", GetFilesByETagDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const ListFilesSummaryDocument = `query ListFilesSummary {
	files {
		... FileMetadataSummaryFragment
	}
}
fragment FileMetadataSummaryFragment on files {
	id
	name
	bucketId
	isUploaded
}
`

func (c *Client) ListFilesSummary(ctx context.Context, interceptors ...clientv2.RequestInterceptor) (*ListFilesSummary, error) {
	vars := map[string]any{}

	var res ListFilesSummary
	if err := c.Client.Post(ctx, "ListFilesSummary", ListFilesSummaryDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const InsertFileDocument = `mutation InsertFile ($object: files_insert_input!) {
	insertFile(object: $object) {
		id
	}
}
`

func (c *Client) InsertFile(ctx context.Context, object FilesInsertInput, interceptors ...clientv2.RequestInterceptor) (*InsertFile, error) {
	vars := map[string]any{
		"object": object,
	}

	var res InsertFile
	if err := c.Client.Post(ctx, "InsertFile", InsertFileDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const UpdateFileDocument = `mutation UpdateFile ($id: uuid!, $_set: files_set_input!) {
	updateFile(pk_columns: {id:$id}, _set: $_set) {
		... FileMetadataFragment
	}
}
fragment FileMetadataFragment on files {
	id
	name
	size
	bucketId
	etag
	createdAt
	updatedAt
	isUploaded
	mimeType
	uploadedByUserId
	metadata
	objectKey
	chunkSize
	chunkCount
	uploadId
}
`

func (c *Client) UpdateFile(ctx context.Context, id string, set FilesSetInput, interceptors ...clientv2.RequestInterceptor) (*UpdateFile, error) {
	vars := map[string]any{
		"id":   id,
		"_set": set,
	}

	var res UpdateFile
	if err := c.Client.Post(ctx, "UpdateFile", UpdateFileDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const DeleteFileDocument = `mutation DeleteFile ($id: uuid!) {
	deleteFile(id: $id) {
		id
	}
}
`

func (c *Client) DeleteFile(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*DeleteFile, error) {
	vars := map[string]any{
		"id": id,
	}

	var res DeleteFile
	if err := c.Client.Post(ctx, "DeleteFile", DeleteFileDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

const InsertVirusDocument = `mutation InsertVirus ($object: virus_insert_input!) {
	insertVirus(object: $object) {
		id
	}
}
`

func (c *Client) InsertVirus(ctx context.Context, object VirusInsertInput, interceptors ...clientv2.RequestInterceptor) (*InsertVirus, error) {
	vars := map[string]any{
		"object": object,
	}

	var res InsertVirus
	if err := c.Client.Post(ctx, "InsertVirus", InsertVirusDocument, &res, vars, interceptors...); err != nil {
		if c.Client.ParseDataWhenErrors {
			return &res, err
		}

		return nil, err
	}

	return &res, nil
}

var DocumentOperationNames = map[string]string{
	GetBucketDocument:        "GetBucket",
	GetFileDocument:          "GetFile",
	GetFilesByETagDocument:   "GetFilesByETag",
	ListFilesSummaryDocument: "ListFilesSummary",
	InsertFileDocument:       "InsertFile",
	UpdateFileDocument:       "UpdateFile",
	DeleteFileDocument:       "DeleteFile",
	InsertVirusDocument:      "InsertVirus",
}
