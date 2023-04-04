package api

import (
	"encoding/json"
	"fmt"
	"marble/marble-backend/app"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ScenarioAppInterface interface {
	GetScenarios(organizationID string) ([]app.Scenario, error)
	CreateScenario(organizationID string, scenario app.Scenario) (app.Scenario, error)

	GetScenario(organizationID string, scenarioID string) (app.Scenario, error)
}

func (a *API) handleScenariosGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, err := orgIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		scenarios, err := a.app.GetScenarios(orgID)
		if err != nil {
			// Could not execute request
			http.Error(w, fmt.Errorf("error getting scenarios: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(&scenarios)
		if err != nil {
			// Could not encode JSON
			http.Error(w, fmt.Errorf("could not encode response JSON: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (a *API) handleScenariosPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, err := orgIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		requestData := &app.Scenario{}
		err = json.NewDecoder(r.Body).Decode(requestData)
		if err != nil {
			// Could not parse JSON
			http.Error(w, fmt.Errorf("could not parse input JSON: %w", err).Error(), http.StatusBadRequest)
			return
		}

		scenario, err := a.app.CreateScenario(orgID, *requestData)
		if err != nil {
			// Could not execute request
			// TODO(errors): handle missing fields error ?
			http.Error(w, fmt.Errorf("error getting scenarios: %w", err).Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(scenario)
		if err != nil {
			// Could not encode JSON
			http.Error(w, fmt.Errorf("could not encode response JSON: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (a *API) handleScenarioGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, err := orgIDFromCtx(r.Context())
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		scenarioID := chi.URLParam(r, "scenarioID")

		scenario, err := a.app.GetScenario(orgID, scenarioID)
		if err != nil {
			// Could not execute request
			http.Error(w, fmt.Errorf("error getting scenario(id: %s): %w", scenarioID, err).Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(&scenario)
		if err != nil {
			// Could not encode JSON
			http.Error(w, fmt.Errorf("could not encode response JSON: %w", err).Error(), http.StatusInternalServerError)
			return
		}
	}
}
