package models

import "time"

type ScheduledExecution struct {
	Id                       string
	OrganizationId           string
	ScenarioId               string
	ScenarioIterationId      string
	Status                   ScheduledExecutionStatus
	StartedAt                time.Time
	FinishedAt               *time.Time
	NumberOfCreatedDecisions int
	Scenario                 Scenario
	Manual                   bool
}

type ScheduledExecutionStatus int

const (
	ScheduledExecutionPending ScheduledExecutionStatus = iota
	ScheduledExecutionProcessing
	ScheduledExecutionSuccess
	ScheduledExecutionFailure
)

func (s ScheduledExecutionStatus) String() string {
	switch s {
	case ScheduledExecutionPending:
		return "pending"
	case ScheduledExecutionProcessing:
		return "processing"
	case ScheduledExecutionSuccess:
		return "success"
	case ScheduledExecutionFailure:
		return "failure"
	}
	return "pending"
}

func ScheduledExecutionStatusFrom(s string) ScheduledExecutionStatus {
	switch s {
	case "pending":
		return ScheduledExecutionPending
	case "success":
		return ScheduledExecutionSuccess
	case "failure":
		return ScheduledExecutionFailure
	case "processing":
		return ScheduledExecutionProcessing
	}
	return ScheduledExecutionPending
}

type UpdateScheduledExecutionStatusInput struct {
	Id                       string
	Status                   ScheduledExecutionStatus
	NumberOfCreatedDecisions *int
	CurrentStatusCondition   ScheduledExecutionStatus // Used for optimistic locking
}

type CreateScheduledExecutionInput struct {
	OrganizationId      string
	ScenarioId          string
	ScenarioIterationId string
	Manual              bool
}

type ListScheduledExecutionsFilters struct {
	OrganizationId string
	ScenarioId     string
	Status         []ScheduledExecutionStatus
	ExcludeManual  bool
}
