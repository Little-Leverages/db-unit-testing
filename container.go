package users

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateContainer(databaseName string) (testcontainers.Container, *sql.DB, error) {
	port := "5432/tcp"

	var env = map[string]string{
		"POSTGRES_USER":     databaseName,
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       databaseName,
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:12.5",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			Name:         databaseName + uuid.New().String(),
			WaitingFor: wait.ForSQL(nat.Port(port), "pgx", func(host string, port nat.Port) string {
				return fmt.Sprintf("postgres://%s:password@%s:%s/%s?sslmode=disable", databaseName, host, port.Port(), databaseName)
			}),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), nat.Port(port))
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s statement_cache_mode=describe", "localhost", mappedPort.Port(), databaseName, "password", databaseName)
	db, err := sql.Open("pgx", url)
	if err != nil {
		fmt.Fprintln(os.Stdout, "stuff here")
		return nil, nil, err
	}

	return container, db, nil
}
