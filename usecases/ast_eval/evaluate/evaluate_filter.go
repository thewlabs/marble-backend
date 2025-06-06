package evaluate

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/cockroachdb/errors"

	"github.com/checkmarble/marble-backend/models"
	"github.com/checkmarble/marble-backend/models/ast"
)

type FilterEvaluator struct {
	DataModel models.DataModel
}

var validTypeForFilterOperators = map[ast.FilterOperator][]models.DataType{
	ast.FILTER_EQUAL:            {models.Bool, models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_NOT_EQUAL:        {models.Bool, models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_GREATER:          {models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_GREATER_OR_EQUAL: {models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_LESSER:           {models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_LESSER_OR_EQUAL:  {models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_IS_IN_LIST:       {models.String},
	ast.FILTER_IS_NOT_IN_LIST:   {models.String},
	ast.FILTER_IS_EMPTY:         {models.Bool, models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_IS_NOT_EMPTY:     {models.Bool, models.Int, models.Float, models.String, models.Timestamp},
	ast.FILTER_STARTS_WITH:      {models.String},
	ast.FILTER_ENDS_WITH:        {models.String},
	ast.FILTER_FUZZY_MATCH:      {models.String},
}

func (f FilterEvaluator) Evaluate(ctx context.Context, arguments ast.Arguments) (any, []error) {
	tableName, tableNameErr := AdaptNamedArgument(arguments.NamedArgs, "tableName", adaptArgumentToString)
	fieldName, fieldNameErr := AdaptNamedArgument(arguments.NamedArgs, "fieldName", adaptArgumentToString)
	operatorStr, operatorErr := AdaptNamedArgument(arguments.NamedArgs, "operator", adaptArgumentToString)
	errs := filterNilErrors(tableNameErr, fieldNameErr, operatorErr)
	if len(errs) > 0 {
		return nil, errs
	}

	fieldType, err := getFieldType(f.DataModel, tableName, fieldName)
	if err != nil {
		return MakeEvaluateError(errors.Join(
			errors.Wrap(err, fmt.Sprintf("field type for %s.%s not found in data model in Evaluate filter", tableName, fieldName)),
			ast.NewNamedArgumentError("fieldName"),
		))
	}

	// Operator validation
	operator := ast.FilterOperator(operatorStr)
	validTypes, isValid := validTypeForFilterOperators[operator]
	if !isValid {
		return MakeEvaluateError(errors.Join(
			errors.Wrap(ast.ErrRuntimeExpression,
				fmt.Sprintf("operator %s is not valid in Evaluate filter", operator)),
			ast.NewNamedArgumentError("operator"),
		))
	}

	isValidFieldType := slices.Contains(validTypes, fieldType)
	if !isValidFieldType {
		return MakeEvaluateError(errors.Join(
			errors.Wrap(ast.ErrArgumentInvalidType,
				fmt.Sprintf("field type %s is not valid for operator %s in Evaluate filter", fieldType.String(), operator)),
			ast.NewNamedArgumentError("fieldName"),
		))
	}

	// Value validation
	value := arguments.NamedArgs["value"]
	if value == nil {
		return ast.Filter{
			TableName: tableName,
			FieldName: fieldName,
			Operator:  operator,
			Value:     nil,
		}, nil
	}

	// The value that is promoted here is then passed directly to the ingested data read repository to be used as a filter value in the sql query.
	var promotedValue any
	switch {
	case operator == ast.FILTER_FUZZY_MATCH:
		// fuzzy match filter takes a custom type (ast.FuzzyMatchOptions), pass it through as it is.
		promotedValue = value
	case fieldType == models.Int && reflect.TypeOf(value) == reflect.TypeOf(float64(0)):
		// When value is a float, it cannot be cast to int but SQL can handle the comparision, so no casting is required
		promotedValue = value
	case operator == ast.FILTER_IS_IN_LIST || operator == ast.FILTER_IS_NOT_IN_LIST:
		// isInList filter takes a slice of strings, accept a slice of any and cast it to a slice of strings (and normalize them)
		promotedValue, err = adaptArgumentToListOfStrings(value)
		// err is checked outside of the switch
	default:
		// by default, it is assumed that the value on the right of the comparison should have the same type as the field on the left
		promotedValue, err = promoteArgumentToDataType(value, fieldType)
		// err is checked outside of the switch
	}
	if err != nil {
		return MakeEvaluateError(errors.Join(
			errors.Wrap(ast.ErrArgumentInvalidType,
				fmt.Sprintf("value is not compatible with selected field %s.%s in Evaluate filter", tableName, fieldName)),
			ast.NewNamedArgumentError("value"),
			err,
		))
	}

	return ast.Filter{
		TableName: tableName,
		FieldName: fieldName,
		Operator:  operator,
		Value:     promotedValue,
	}, nil
}
