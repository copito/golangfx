package stream_profiling

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/copito/runner/src/pkg/utils"
	"github.com/copito/runner/src/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationPostgresExample(t *testing.T) {
	// skip integration test in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	rootPath, err := utils.GetModuleRoot()
	if err != nil {
		t.Fatalf("failed to get module root: %v", err)
	}

	ctx := context.Background()
	testContainer := setup.NewPostgres(ctx, setup.PostgresConfig{
		User:            "postgres",
		Password:        "postgres",
		DBName:          "testdb",
		InitScriptsPath: []string{filepath.Join(rootPath, "src", "tests", "testdata", "sql", "postgres", "postgres-simple-example.sql")},
		ConfigFilePath:  filepath.Join(rootPath, "src", "tests", "testdata", "configs", "postgres.conf"),
	})
	defer testContainer.Teardown(ctx)

	rc, _ := testContainer.Container.Logs(ctx)
	defer rc.Close()
	t.Log(rc)

	connectionString, err := testContainer.Container.ConnectionString(ctx, "sslmode=disable", "application_name=test")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	t.Logf("Postgres connection string: %s", connectionString)

	// time.Sleep(1000 * time.Second)

	conn, err := sql.Open("pgx", connectionString)
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer conn.Close()

	err = conn.PingContext(ctx)
	if err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	t.Run("testing that we have 3 rows", func(t *testing.T) {
		query := `SELECT COUNT(*) FROM testdb.example_schema.example_table`
		var value int
		err = conn.QueryRowContext(ctx, query).Scan(&value)
		if err != nil {
			t.Fatalf("failed to execute query: %v", err)
		}

		// based on testdata/sql/postgres
		assert.Equal(t, 3, value, "expected a row count to be 3")
	})
}
