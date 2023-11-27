package dbmodels

import (
	"github.com/checkmarble/marble-backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBCase struct {
	Id             pgtype.Text      `db:"id"`
	CreatedAt      pgtype.Timestamp `db:"created_at"`
	InboxId        pgtype.Text      `db:"inbox_id"`
	Name           pgtype.Text      `db:"name"`
	OrganizationId pgtype.Text      `db:"org_id"`
	Status         pgtype.Text      `db:"status"`
}

type DBCaseWithContributors struct {
	DBCase
	Contributors   []DBCaseContributor `db:"contributors"`
	DecisionsCount int                 `db:"decisions_count"`
}

const TABLE_CASES = "cases"

var SelectCaseColumn = []string{"id", "created_at", "inbox_id", "name", "org_id", "status"}

func AdaptCase(db DBCase) (models.Case, error) {
	return models.Case{
		Id:             db.Id.String,
		CreatedAt:      db.CreatedAt.Time,
		InboxId:        db.InboxId.String,
		Name:           db.Name.String,
		OrganizationId: db.OrganizationId.String,
		Status:         models.CaseStatus(db.Status.String),
	}, nil
}

func AdaptCaseWithContributors(db DBCaseWithContributors) (models.Case, error) {
	caseModel, err := AdaptCase(db.DBCase)
	if err != nil {
		return models.Case{}, err
	}
	caseModel.DecisionsCount = db.DecisionsCount

	caseModel.Contributors = make([]models.CaseContributor, len(db.Contributors))
	for i, contributor := range db.Contributors {
		caseModel.Contributors[i], err = AdaptCaseContributor(contributor)
		if err != nil {
			return models.Case{}, err
		}
	}

	return caseModel, nil
}
