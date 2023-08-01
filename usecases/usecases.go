package usecases

import (
	"marble/marble-backend/models"
	"marble/marble-backend/models/ast"
	"marble/marble-backend/repositories"
	"marble/marble-backend/usecases/ast_eval"
	"marble/marble-backend/usecases/ast_eval/evaluate"
	"marble/marble-backend/usecases/org_transaction"
	"marble/marble-backend/usecases/organization"
	"marble/marble-backend/usecases/scenarios"
	"marble/marble-backend/usecases/scheduledexecution"
	"marble/marble-backend/usecases/security"
)

type Usecases struct {
	Repositories  repositories.Repositories
	Configuration models.GlobalConfiguration
}

func (usecases *Usecases) NewOrganizationUseCase() OrganizationUseCase {
	return OrganizationUseCase{
		transactionFactory:           usecases.Repositories.TransactionFactory,
		orgTransactionFactory:        usecases.NewOrgTransactionFactory(),
		organizationRepository:       usecases.Repositories.OrganizationRepository,
		datamodelRepository:          usecases.Repositories.DataModelRepository,
		apiKeyRepository:             usecases.Repositories.ApiKeyRepository,
		userRepository:               usecases.Repositories.UserRepository,
		organizationCreator:          usecases.NewOrganizationCreator(),
		organizationSchemaRepository: usecases.Repositories.OrganizationSchemaRepository,
		populateOrganizationSchema:   usecases.NewPopulateOrganizationSchema(),
	}
}

func (usecases *Usecases) NewOrgTransactionFactory() org_transaction.Factory {
	return &org_transaction.FactoryImpl{
		OrganizationSchemaRepository:     usecases.Repositories.OrganizationSchemaRepository,
		TransactionFactory:               usecases.Repositories.TransactionFactory,
		DatabaseConnectionPoolRepository: usecases.Repositories.DatabaseConnectionPoolRepository,
	}
}

func (usecases *Usecases) NewIngestionUseCase() IngestionUseCase {
	return IngestionUseCase{
		orgTransactionFactory: usecases.NewOrgTransactionFactory(),
		ingestionRepository:   usecases.Repositories.IngestionRepository,
		gcsRepository:         usecases.Repositories.GcsRepository,
		datamodelRepository:   usecases.Repositories.DataModelRepository,
	}
}

func (usecases *Usecases) NewDecisionUsecase() DecisionUsecase {
	return DecisionUsecase{
		transactionFactory:              usecases.Repositories.TransactionFactory,
		orgTransactionFactory:           usecases.NewOrgTransactionFactory(),
		ingestedDataReadRepository:      usecases.Repositories.IngestedDataReadRepository,
		decisionRepository:              usecases.Repositories.DecisionRepository,
		datamodelRepository:             usecases.Repositories.DataModelRepository,
		scenarioReadRepository:          usecases.Repositories.ScenarioReadRepository,
		scenarioIterationReadRepository: usecases.Repositories.ScenarioIterationReadRepository,
		customListRepository:            usecases.Repositories.CustomListRepository,
		evaluateRuleAstExpression:       usecases.NewEvaluateRuleAstExpression(),
	}
}

func (usecases *Usecases) NewUserUseCase() UserUseCase {
	return UserUseCase{
		transactionFactory: usecases.Repositories.TransactionFactory,
		userRepository:     usecases.Repositories.UserRepository,
	}
}

func (usecases *Usecases) NewCustomListUseCase() CustomListUseCase {
	return CustomListUseCase{
		transactionFactory:   usecases.Repositories.TransactionFactory,
		CustomListRepository: usecases.Repositories.CustomListRepository,
	}
}

func (usecases *Usecases) NewSeedUseCase() SeedUseCase {
	return SeedUseCase{
		transactionFactory:     usecases.Repositories.TransactionFactory,
		userRepository:         usecases.Repositories.UserRepository,
		organizationCreator:    usecases.NewOrganizationCreator(),
		organizationRepository: usecases.Repositories.OrganizationRepository,
		customListRepository:   usecases.Repositories.CustomListRepository,
	}
}

func (usecases *Usecases) NewOrganizationCreator() organization.OrganizationCreator {
	return organization.OrganizationCreator{
		TransactionFactory:     usecases.Repositories.TransactionFactory,
		OrganizationRepository: usecases.Repositories.OrganizationRepository,
		DataModelRepository:    usecases.Repositories.DataModelRepository,
		OrganizationSeeder: &organization.OrganizationSeederImpl{
			CustomListRepository:             usecases.Repositories.CustomListRepository,
			ApiKeyRepository:                 usecases.Repositories.ApiKeyRepository,
			ScenarioWriteRepository:          usecases.Repositories.ScenarioWriteRepository,
			ScenarioIterationWriteRepository: usecases.Repositories.ScenarioIterationWriteRepository,
			ScenarioPublisher:                usecases.NewScenarioPublisher(),
		},
		PopulateOrganizationSchema: usecases.NewPopulateOrganizationSchema(),
	}
}

func (usecases *Usecases) NewScenarioIterationRuleUsecase() ScenarioIterationRuleUsecase {
	return ScenarioIterationRuleUsecase{
		repository:                usecases.Repositories.ScenarioIterationRuleRepositoryLegacy,
		scenarioFetcher:           usecases.NewScenarioFetcher(),
		validateScenarioIteration: usecases.NewValidateScenarioIteration(),
	}
}

func (usecases *Usecases) NewScheduledExecutionUsecase() ScheduledExecutionUsecase {
	return ScheduledExecutionUsecase{
		scenarioReadRepository:          usecases.Repositories.ScenarioReadRepository,
		scenarioIterationReadRepository: usecases.Repositories.ScenarioIterationReadRepository,
		scheduledExecutionRepository:    usecases.Repositories.ScheduledExecutionRepository,
		transactionFactory:              usecases.Repositories.TransactionFactory,
		orgTransactionFactory:           usecases.NewOrgTransactionFactory(),
		dataModelRepository:             usecases.Repositories.DataModelRepository,
		ingestedDataReadRepository:      usecases.Repositories.IngestedDataReadRepository,
		decisionRepository:              usecases.Repositories.DecisionRepository,
		scenarioPublicationsRepository:  usecases.Repositories.ScenarioPublicationRepository,
		customListRepository:            usecases.Repositories.CustomListRepository,
		exportScheduleExecution:         usecases.NewExportScheduleExecution(),
		evaluateRuleAstExpression:       usecases.NewEvaluateRuleAstExpression(),
	}
}

func (usecases *Usecases) NewExportScheduleExecution() scheduledexecution.ExportScheduleExecution {
	return &scheduledexecution.ExportScheduleExecutionImpl{
		AwsS3Repository:        usecases.Repositories.AwsS3Repository,
		DecisionRepository:     usecases.Repositories.DecisionRepository,
		OrganizationRepository: usecases.Repositories.OrganizationRepository,
	}
}

func (usecases *Usecases) NewPopulateOrganizationSchema() organization.PopulateOrganizationSchema {
	return organization.PopulateOrganizationSchema{
		TransactionFactory:           usecases.Repositories.TransactionFactory,
		OrganizationRepository:       usecases.Repositories.OrganizationRepository,
		OrganizationSchemaRepository: usecases.Repositories.OrganizationSchemaRepository,
		DataModelRepository:          usecases.Repositories.DataModelRepository,
	}
}

func (usecases *Usecases) AstEvaluationEnvironment(organizationId string, payload models.PayloadReader) ast_eval.AstEvaluationEnvironment {
	environment := ast_eval.NewAstEvaluationEnvironment()

	// execution of a scenario with a dedicated security context
	enforceSecurity := &security.EnforceSecurityImpl{
		Credentials: models.Credentials{
			OrganizationId: organizationId,
		},
	}

	environment.AddEvaluator(ast.FUNC_CUSTOM_LIST_ACCESS,
		evaluate.NewCustomListValuesAccess(
			usecases.Repositories.CustomListRepository,
			enforceSecurity,
		),
	)

	environment.AddEvaluator(ast.FUNC_DB_ACCESS,
		evaluate.NewDatabaseAccess(
			usecases.NewOrgTransactionFactory(),
			usecases.Repositories.IngestedDataReadRepository,
			usecases.Repositories.DataModelRepository,
			payload,
			organizationId,
		))
	environment.AddEvaluator(ast.FUNC_PAYLOAD, evaluate.NewPayload(ast.FUNC_PAYLOAD, payload))
	return environment
}

func (usecases *Usecases) NewEvaluateRuleAstExpression() ast_eval.EvaluateRuleAstExpression {
	return ast_eval.EvaluateRuleAstExpression{
		AstEvaluationEnvironmentFactory: usecases.AstEvaluationEnvironment,
	}
}

func (usecases *Usecases) NewScenarioPublisher() scenarios.ScenarioPublisher {
	return scenarios.ScenarioPublisher{
		ScenarioPublicationsRepository:  usecases.Repositories.ScenarioPublicationRepository,
		ScenarioWriteRepository:         usecases.Repositories.ScenarioWriteRepository,
		ScenarioIterationReadRepository: usecases.Repositories.ScenarioIterationReadRepository,
		ValidateScenarioIteration:       usecases.NewValidateScenarioIteration(),
	}
}

func (usecases *Usecases) NewValidateScenarioIteration() scenarios.ValidateScenarioIteration {
	return &scenarios.ValidateScenarioIterationImpl{
		DataModelRepository:             usecases.Repositories.DataModelRepository,
		AstEvaluationEnvironmentFactory: usecases.AstEvaluationEnvironment,
	}
}

func (usecase *Usecases) NewScenarioFetcher() scenarios.ScenarioFetcher {
	return scenarios.ScenarioFetcher{
		ScenarioReadRepository:          usecase.Repositories.ScenarioReadRepository,
		ScenarioIterationReadRepository: usecase.Repositories.ScenarioIterationReadRepository,
	}
}
