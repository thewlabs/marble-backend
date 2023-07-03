package evaluate

import (
	"context"
	"errors"
	"fmt"
	"marble/marble-backend/models"
	"marble/marble-backend/models/ast"
	"marble/marble-backend/repositories"
	"marble/marble-backend/usecases/organization"
)

type DatabaseAccess struct {
	DatamodelRepository        repositories.DataModelRepository
	Payload                    models.PayloadReader
	OrgTransactionFactory      organization.OrgTransactionFactory
	IngestedDataReadRepository repositories.IngestedDataReadRepository
	Creds                      models.Credentials
}

func NewDatabaseAccess(otf organization.OrgTransactionFactory, idrr repositories.IngestedDataReadRepository,
	dm repositories.DataModelRepository, payload models.PayloadReader,
	creds models.Credentials) DatabaseAccess {
	return DatabaseAccess{
		DatamodelRepository:        dm,
		Payload:                    payload,
		OrgTransactionFactory:      otf,
		IngestedDataReadRepository: idrr,
		Creds:                      creds,
	}
}

func (d DatabaseAccess) Evaluate(arguments ast.Arguments) (any, error) {
	var pathStringArr []string
	tableName, ok := arguments.NamedArgs["tableName"].(string)
	if !ok {
		return nil, fmt.Errorf("tableName is not a string %w", ErrRuntimeExpression)
	}
	fieldName, ok := arguments.NamedArgs["fieldName"].(string)
	if !ok {
		return nil, fmt.Errorf("fieldName is not a string %w", ErrRuntimeExpression)
	}
	path, ok := arguments.NamedArgs["path"].([]any)
	if !ok {
		return nil, fmt.Errorf("path is not a string %w", ErrRuntimeExpression)
	}
	for _, v := range path {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("path value is not a string %w", ErrRuntimeExpression)
		}
		pathStringArr = append(pathStringArr, str)
	}

	fieldValue, err := d.getDbField(tableName, fieldName, pathStringArr)
	if err != nil {
		return nil, fmt.Errorf("getDbField not working %w", err)
	}
	return fieldValue, nil
}

func (d DatabaseAccess) getDbField(tableName string, fieldName string, path []string) (interface{}, error) {
	dm, err := d.DatamodelRepository.GetDataModel(context.Background(), d.Creds.OrganizationId)
	if errors.Is(err, models.NotFoundInRepositoryError) {
		return models.Decision{}, fmt.Errorf("data model not found: %w", models.NotFoundError)
	} else if err != nil {
		return models.Decision{}, fmt.Errorf("error getting data model: %w", err)
	}

	return organization.TransactionInOrgSchemaReturnValue(
		d.OrgTransactionFactory,
		d.Creds.OrganizationId,
		func(tx repositories.Transaction) (interface{}, error) {
			return d.IngestedDataReadRepository.GetDbField(tx, models.DbFieldReadParams{
				TriggerTableName: models.TableName(tableName),
				Path:             models.ToLinkNames(path),
				FieldName:        models.FieldName(fieldName),
				DataModel:        dm,
				Payload:          d.Payload,
			})
		})
}
