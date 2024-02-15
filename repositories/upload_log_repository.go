package repositories

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/checkmarble/marble-backend/models"
	"github.com/checkmarble/marble-backend/repositories/dbmodels"
)

type UploadLogRepository interface {
	CreateUploadLog(ctx context.Context, exec Executor, log models.UploadLog) error
	UpdateUploadLog(ctx context.Context, exec Executor, input models.UpdateUploadLogInput) error
	UploadLogById(ctx context.Context, exec Executor, id string) (models.UploadLog, error)
	AllUploadLogsByStatus(ctx context.Context, exec Executor, status models.UploadStatus) ([]models.UploadLog, error)
	AllUploadLogsByTable(ctx context.Context, exec Executor, organizationId, tableName string) ([]models.UploadLog, error)
}

type UploadLogRepositoryImpl struct{}

func (repo *UploadLogRepositoryImpl) CreateUploadLog(ctx context.Context, exec Executor, log models.UploadLog) error {
	if err := validateMarbleDbExecutor(exec); err != nil {
		return err
	}

	err := ExecBuilder(
		ctx,
		exec,
		NewQueryBuilder().Insert(dbmodels.TABLE_UPLOAD_LOGS).
			Columns(
				"id",
				"org_id",
				"user_id",
				"file_name",
				"table_name",
				"status",
				"started_at",
				"finished_at",
				"lines_processed",
			).
			Values(
				log.Id,
				log.OrganizationId,
				log.UserId,
				log.FileName,
				log.TableName,
				log.UploadStatus,
				log.StartedAt,
				log.FinishedAt,
				log.LinesProcessed,
			),
	)
	return err
}

func (repo *UploadLogRepositoryImpl) UpdateUploadLog(ctx context.Context, exec Executor, input models.UpdateUploadLogInput) error {
	if err := validateMarbleDbExecutor(exec); err != nil {
		return err
	}

	var updateRequest = NewQueryBuilder().Update(dbmodels.TABLE_UPLOAD_LOGS)

	if input.UploadStatus != "" {
		updateRequest = updateRequest.Set("status", input.UploadStatus)
	}
	if input.FinishedAt != nil {
		updateRequest = updateRequest.Set("finished_at", *input.FinishedAt)
	}
	updateRequest = updateRequest.Where("id = ?", input.Id)

	err := ExecBuilder(ctx, exec, updateRequest)
	return err
}

func (repo *UploadLogRepositoryImpl) UploadLogById(ctx context.Context, exec Executor, id string) (models.UploadLog, error) {
	if err := validateMarbleDbExecutor(exec); err != nil {
		return models.UploadLog{}, err
	}

	uploadLog, err := SqlToModel(
		ctx,
		exec,
		NewQueryBuilder().
			Select(dbmodels.SelectUploadLogColumn...).
			From(dbmodels.TABLE_UPLOAD_LOGS).
			Where(squirrel.Eq{"id": id}),
		dbmodels.AdaptUploadLog,
	)

	if err != nil {
		return models.UploadLog{}, err
	}

	return uploadLog, err
}

func (repo *UploadLogRepositoryImpl) AllUploadLogsByStatus(ctx context.Context, exec Executor, status models.UploadStatus) ([]models.UploadLog, error) {
	if err := validateMarbleDbExecutor(exec); err != nil {
		return nil, err
	}

	return SqlToListOfModels(
		ctx,
		exec,
		NewQueryBuilder().
			Select(dbmodels.SelectUploadLogColumn...).
			From(dbmodels.TABLE_UPLOAD_LOGS).
			Where(squirrel.Eq{"status": status}).
			OrderBy("started_at"),
		dbmodels.AdaptUploadLog,
	)
}

func (repo *UploadLogRepositoryImpl) AllUploadLogsByTable(ctx context.Context, exec Executor, organizationId, tableName string) ([]models.UploadLog, error) {
	if err := validateMarbleDbExecutor(exec); err != nil {
		return nil, err
	}

	return SqlToListOfModels(
		ctx,
		exec,
		NewQueryBuilder().
			Select(dbmodels.SelectUploadLogColumn...).
			From(dbmodels.TABLE_UPLOAD_LOGS).
			Where(squirrel.Eq{"org_id": organizationId}).
			Where(squirrel.Eq{"table_name": tableName}).
			OrderBy("started_at DESC"),
		dbmodels.AdaptUploadLog,
	)
}
