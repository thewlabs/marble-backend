package repositories

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/checkmarble/marble-backend/models"
	"github.com/checkmarble/marble-backend/utils"
	"github.com/cockroachdb/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

const signedUrlExpiryHours = 1

type GcsRepository interface {
	ListFiles(ctx context.Context, bucketName, prefix string) ([]models.GCSFile, error)
	GetFile(ctx context.Context, bucketName, fileName string) (models.GCSFile, error)
	MoveFile(ctx context.Context, bucketName, source, destination string) error
	OpenStream(ctx context.Context, bucketName, fileName string) io.WriteCloser
	DeleteFile(ctx context.Context, bucketName, fileName string) error
	UpdateFileMetadata(ctx context.Context, bucketName, fileName string, metadata map[string]string) error
	GenerateSignedUrl(ctx context.Context, bucketName, fileName string) (string, error)
}

type GcsRepositoryImpl struct {
	gcsClient *storage.Client
}

func (repository *GcsRepositoryImpl) getGCSClient(ctx context.Context) *storage.Client {
	// Lazy load the GCS client, as it is used only in one batch usecase, to avoid requiring GCS credentials for all devs
	if repository.gcsClient != nil {
		return repository.gcsClient
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to load GCS client: %w", err))
	}
	repository.gcsClient = client
	return client
}

// Not used since legacy CSV ingestion has been removed
func (repository *GcsRepositoryImpl) ListFiles(ctx context.Context, bucketName, prefix string) ([]models.GCSFile, error) {
	bucket := repository.getGCSClient(ctx).Bucket(bucketName)
	_, err := bucket.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket to list GCS objects from bucket %s/%s: %w", bucketName, prefix, err)
	}

	var output []models.GCSFile

	query := &storage.Query{Prefix: prefix}
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done { //nolint:errorlint
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list GCS objects from bucket %s/%s: %w", bucketName, prefix, err)
		}

		r, err := bucket.Object(attrs.Name).NewReader(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to read GCS object %s/%s: %w", bucketName, attrs.Name, err)
		}

		output = append(output, models.GCSFile{
			FileName:   attrs.Name,
			Reader:     r,
			BucketName: bucketName,
		})
	}

	return output, nil
}

func (repository *GcsRepositoryImpl) GetFile(ctx context.Context, bucketName, fileName string) (models.GCSFile, error) {
	tracer := utils.OpenTelemetryTracerFromContext(ctx)
	ctx, span := tracer.Start(
		ctx,
		"repositories.GcsRepository.GetFile",
		trace.WithAttributes(attribute.String("bucket", bucketName)),
		trace.WithAttributes(attribute.String("fileName", fileName)),
	)
	defer span.End()
	bucket := repository.getGCSClient(ctx).Bucket(bucketName)

	ctxBucket, span2 := tracer.Start(
		ctx,
		"repositories.GcsRepository.GetFile - bucket attrs",
	)
	defer span2.End()
	_, err := bucket.Attrs(ctxBucket)
	if err != nil {
		return models.GCSFile{}, fmt.Errorf("failed to get bucket %s: %w", bucketName, err)
	}
	span2.End()

	ctx, span = tracer.Start(
		ctx,
		"repositories.GcsRepository.GetFile - file reader",
	)
	defer span.End()
	reader, err := bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		return models.GCSFile{}, fmt.Errorf("failed to read GCS object %s/%s: %w", bucketName, fileName, err)
	}

	return models.GCSFile{
		FileName:   fileName,
		Reader:     reader,
		BucketName: bucketName,
	}, nil
}

// Not used since legacy CSV ingestion has been removed
func (repository *GcsRepositoryImpl) MoveFile(ctx context.Context, bucketName, srcName, destName string) error {
	gcsClient := repository.getGCSClient(ctx)
	src := gcsClient.Bucket(bucketName).Object(srcName)
	dst := gcsClient.Bucket(bucketName).Object(destName)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy the file is aborted
	// if the object's generation number does not match your precondition.
	// For a dst object that does not yet exist, set the DoesNotExist precondition.
	// Straight from the docs: https://cloud.google.com/storage/docs/copying-renaming-moving-objects?hl=fr#move
	dst = dst.If(storage.Conditions{DoesNotExist: true})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %w", destName, srcName, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", srcName, err)
	}
	return nil
}

func (repository *GcsRepositoryImpl) OpenStream(ctx context.Context, bucketName, fileName string) io.WriteCloser {
	gcsClient := repository.getGCSClient(ctx)

	writer := gcsClient.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	writer.ChunkSize = 0 // note retries are not supported for chunk size 0.
	return writer
}

func (repository *GcsRepositoryImpl) UpdateFileMetadata(ctx context.Context,
	bucketName, fileName string, metadata map[string]string,
) error {
	gcsClient := repository.getGCSClient(ctx)
	defer gcsClient.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	object := gcsClient.Bucket(bucketName).Object(fileName)

	// Optional: set a metageneration-match precondition to avoid potential race
	// conditions and data corruptions. The request to update is aborted if the
	// object's metageneration does not match your precondition.
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	object = object.If(storage.Conditions{MetagenerationMatch: attrs.Metageneration})

	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{Metadata: metadata}

	if _, err := object.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %w", fileName, err)
	}

	return nil
}

func (repository *GcsRepositoryImpl) DeleteFile(ctx context.Context, bucketName, fileName string) error {
	gcsClient := repository.getGCSClient(ctx)
	defer gcsClient.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	object := gcsClient.Bucket(bucketName).Object(fileName)

	if err := object.Delete(ctx); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error deleting file: %s", fileName))
	}

	return nil
}

func (repo *GcsRepositoryImpl) GenerateSignedUrl(ctx context.Context, bucketName, fileName string) (string, error) {
	// This code will typically not run locally if you target the real GCS repository, because SignedURL only works with service account credentials (not end user credentials)
	// Hence, run the code locally with the fake GCS repository always
	bucket := repo.getGCSClient(ctx).Bucket(bucketName)
	return bucket.
		SignedURL(
			fileName,
			&storage.SignedURLOptions{
				Method:  http.MethodGet,
				Expires: time.Now().Add(signedUrlExpiryHours * time.Hour),
			},
		)
}
