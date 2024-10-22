package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresContainerName = "movies-reviews-it-postgres"
	postgresUser          = "user"
	postgresPassword      = "pass"
	postgresDb            = "moviesreviewsdb"
)

func prepareInfrastructure(t *testing.T, runFunc func(t *testing.T, connString string)) {
	// Start Postgres container
	postgres, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:  postgresContainerName,
			Image: "postgres:17-alpine",
			Env: map[string]string{
				"POSTGRES_USER":     postgresUser,
				"POSTGRES_PASSWORD": postgresPassword,
				"POSTGRES_DB":       postgresDb,
			},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer cleanUp(t, postgres.Terminate)

	postgresPort, err := postgres.MappedPort(context.Background(), "5432")
	require.NoError(t, err)
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", postgresUser, postgresPassword, "localhost", postgresPort.Int(), postgresDb)

	// Run migrations
	time.Sleep(time.Second * 2)
	runMigrations(t, pgConnString)

	// Run integration_tests
	runFunc(t, pgConnString)
}

func runMigrations(t *testing.T, connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	require.NoError(t, err)

	migrator, err := migrate.NewMigrator(context.Background(), conn, "migrations")
	require.NoError(t, err)

	err = migrator.LoadMigrations("../tern/migrations")
	require.NoError(t, err)

	/*
		migrator.Migrations = append(migrator.Migrations, &migrate.Migration{
			UpSQL: `CREATE TYPE role AS ENUM ('admin', 'editor', 'user');
							CREATE TABLE users (
			    				id SERIAL PRIMARY KEY,
			    				username VARCHAR(24) UNIQUE NOT NULL,
			    				email VARCHAR(128) UNIQUE NOT NULL,
			    				pass_hash VARCHAR(60) NOT NULL,
			    				role role NOT NULL DEFAULT 'user',
			    				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			    				deleted_at TIMESTAMP);`,
			DownSQL: `DROP TABLE IF EXISTS users;
							DROP TYPE IF EXISTS role;`,
		})
	*/

	err = migrator.Migrate(context.Background())
	require.NoError(t, err)
}

func cleanUp(t *testing.T, terminate func(ctx context.Context) error) {
	require.NoError(t, terminate(context.Background()))
}
