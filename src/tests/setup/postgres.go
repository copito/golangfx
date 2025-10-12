package setup

import (
	"context"
	"log"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
}

type PostgresConfig struct {
	User            string
	Password        string
	DBName          string
	InitScriptsPath []string
	ConfigFilePath  string
}

func NewPostgres(ctx context.Context, config PostgresConfig) *PostgresContainer {
	container, err := SetupPostgres(ctx, config)
	if err != nil {
		log.Fatalf("failed to setup Postgres container: %v", err)
	}
	return &PostgresContainer{
		Container: container,
	}
}

func SetupPostgres(ctx context.Context, config PostgresConfig) (*postgres.PostgresContainer, error) {
	var initScripts []string
	if len(config.InitScriptsPath) > 0 {
		initScripts = config.InitScriptsPath
	} else {
		initScripts = []string{filepath.Join("testdata", "sql", "postgres", "postgres-init-user-db.sh")}
	}

	configFile := config.ConfigFilePath
	if configFile == "" {
		configFile = filepath.Join("testdata", "configs", "postgres.conf")
	}

	dbName := config.DBName
	if dbName == "" {
		dbName = "public"
	}

	dbUser := config.User
	if dbUser == "" {
		dbUser = "user"
	}

	dbPassword := config.Password
	if dbPassword == "" {
		dbPassword = "password"
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(initScripts...),
		postgres.WithConfigFile(configFile),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	return postgresContainer, nil
}

func (s *PostgresContainer) Teardown(ctx context.Context) error {
	err := testcontainers.TerminateContainer(s.Container)
	return err
}
