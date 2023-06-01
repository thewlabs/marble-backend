package repositories

import (
	"context"
	"marble/marble-backend/app"
)

type ScenarioReadRepository interface {
	GetScenario(ctx context.Context, orgID string, scenarioID string) (app.Scenario, error)
}
