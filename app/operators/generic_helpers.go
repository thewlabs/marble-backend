package operators

import (
	"context"
	"errors"
	"fmt"
	"marble/marble-backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getPayloadFieldGeneric[T string | bool | float64](d DataAccessor, fieldName string) (T, error) {
	var output T
	// Get account_id from payload. Basically a copy/paste from the string payload field operator Eval()
	fieldRaw, err := d.GetPayloadField(fieldName)
	if err != nil {
		return output, err
	}

	fieldPointer, ok := fieldRaw.(*T)
	if !ok {
		return output, fmt.Errorf("Payload field %s is not a pointer to the right type %T", "account_id", output)
	}
	if fieldPointer == nil {
		return output, fmt.Errorf("Payload field %s is null: %w", "account_id", models.OperatorNullValueReadError)
	}
	output = *fieldPointer

	return output, nil
}

func queryDbFieldGeneric[T float64 | string](ctx context.Context, db *pgxpool.Pool, sql string, args []any) (T, error) {
	var output T
	rows := db.QueryRow(ctx, sql, args...)
	err := rows.Scan(&output)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return output, fmt.Errorf("No rows scanned while reading DB: %w", models.OperatorNoRowsReadInDbError)
	}
	return output, err
}
