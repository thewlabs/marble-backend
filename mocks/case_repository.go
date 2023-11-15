package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/checkmarble/marble-backend/models"
	"github.com/checkmarble/marble-backend/repositories"
)

type CaseRepository struct {
	mock.Mock
}

func (r *CaseRepository) ListCases(tx repositories.Transaction, organizationId string, filters models.CaseFilters) ([]models.Case, error) {
	args := r.Called(tx, organizationId)
	return args.Get(0).([]models.Case), args.Error(1)
}

func (r *CaseRepository) GetCaseById(tx repositories.Transaction, caseId string) (models.Case, error) {
	args := r.Called(tx, caseId)
	return args.Get(0).(models.Case), args.Error(1)
}

func (r *CaseRepository) CreateCase(tx repositories.Transaction, createCaseAttributes models.CreateCaseAttributes, newCaseId string) error {
	args := r.Called(tx, createCaseAttributes, newCaseId)
	return args.Error(0)
}
