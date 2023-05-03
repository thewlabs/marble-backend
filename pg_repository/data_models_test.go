package pg_repository

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"marble/marble-backend/app"
)

type dataModelTestCase struct {
	name           string
	input          app.DataModel
	expectedOutput interface{}
}

func TestDataModelRepoEndToEnd(t *testing.T) {
	transactions := app.Table{
		Name: "transactions_test",
		Fields: map[string]app.Field{
			"object_id": {
				DataType: app.String,
			},
			"updated_at":  {DataType: app.Timestamp},
			"value":       {DataType: app.Float},
			"isValidated": {DataType: app.Bool},
			"account_id":  {DataType: app.String},
		},
		LinksToSingle: map[string]app.LinkToSingle{
			"bank_accounts_test": {
				LinkedTableName: "bank_accounts_test",
				ParentFieldName: "object_id",
				ChildFieldName:  "account_id",
			},
		},
	}
	bank_accounts := app.Table{
		Name: "bank_accounts_test",
		Fields: map[string]app.Field{
			"object_id": {
				DataType: app.String,
			},
			"updated_at":   {DataType: app.Timestamp},
			"status":       {DataType: app.String},
			"is_validated": {DataType: app.Bool},
		},
		LinksToSingle: map[string]app.LinkToSingle{},
	}

	dataModel := app.DataModel{
		Tables: map[string]app.Table{
			"transactions_test":  transactions,
			"bank_accounts_test": bank_accounts,
		},
		Version: "1.0.0",
	}
	ctx := context.Background()

	cases := []dataModelTestCase{
		{
			name:           "Read boolean field from DB without join",
			input:          dataModel,
			expectedOutput: dataModel,
		},
	}

	orgID := globalTestParams.testIds["OrganizationId"]
	asserts := assert.New(t)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			val, err := TestRepo.CreateDataModel(ctx, orgID, dataModel)
			if err != nil {
				t.Errorf("Could not read field from DB: %s", err)
			}

			asserts.Equal(c.expectedOutput, val, "[Create] Output data model should match the input one")

			val, err = TestRepo.GetDataModel(ctx, orgID)
			if err != nil {
				t.Errorf("Could not read field from DB: %s", err)
			}
			asserts.Equal(c.expectedOutput, val, "[Get] Output data model should match the input one")

			unknownOrgID, _ := uuid.NewV4()
			val, err = TestRepo.GetDataModel(ctx, unknownOrgID.String())
			if !errors.Is(err, app.ErrNotFoundInRepository) {
				t.Errorf("Should return an error if the org id is unknown: %s", err)
			}
		})

	}

}
