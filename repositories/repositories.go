package repositories

import (
	"crypto/rsa"
	"marble/marble-backend/models"
	"marble/marble-backend/pg_repository"

	"firebase.google.com/go/v4/auth"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	TransactionFactory              TransactionFactory
	FirebaseTokenRepository         FireBaseTokenRepository
	MarbleJwtRepository             MarbleJwtRepository
	UserRepository                  UserRepository
	ApiKeyRepository                ApiKeyRepository
	OrganizationRepository          OrganizationRepository
	IngestionRepository             IngestionRepository
	DataModelRepository             DataModelRepository
	DbPoolRepository                DbPoolRepository
	IngestedDataReadRepository      IngestedDataReadRepository
	DecisionRepository              DecisionRepository
	ScenarioReadRepository          ScenarioReadRepository
	ScenarioIterationReadRepository ScenarioIterationReadRepository
}

func NewRepositories(
	marbleJwtSigningKey rsa.PrivateKey,
	firebaseClient auth.Client,
	pgRepository *pg_repository.PGRepository,
	marbleConnectionPool *pgxpool.Pool,
	clientConnectionStrings map[models.DatabaseName]string,
) *Repositories {
	databaseConnectionPoolRepository := NewDatabaseConnectionPoolRepository(
		marbleConnectionPool,
		clientConnectionStrings,
	)

	transactionFactory := &TransactionFactoryPosgresql{
		DatabaseConnectionPoolRepository: databaseConnectionPoolRepository,
	}

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Repositories{
		TransactionFactory: transactionFactory,
		FirebaseTokenRepository: FireBaseTokenRepository{
			firebaseClient: firebaseClient,
		},
		MarbleJwtRepository: MarbleJwtRepository{
			jwtSigningPrivateKey: marbleJwtSigningKey,
		},
		UserRepository: &UserRepositoryPostgresql{
			queryBuilder: queryBuilder,
		},
		ApiKeyRepository:                pgRepository,
		OrganizationRepository:          pgRepository,
		IngestionRepository:             pgRepository,
		DataModelRepository:             pgRepository,
		DbPoolRepository:                pgRepository,
		IngestedDataReadRepository:      pgRepository,
		DecisionRepository:              pgRepository,
		ScenarioReadRepository:          pgRepository,
		ScenarioIterationReadRepository: pgRepository,
	}
}
