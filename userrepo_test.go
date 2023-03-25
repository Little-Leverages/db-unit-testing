package users

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

var testSqlDB *sql.DB

func TestMain(m *testing.M) {
	container, testDB, err := CreateContainer("test-db")
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	defer container.Terminate(context.Background()) // nolint

	testSqlDB = testDB

	// Seed data if needed here
	content, err := os.ReadFile("setup.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = testDB.Exec(string(content))
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestCreateAndGet(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		userId := uuid.NewString()
		user := User{
			ID:   userId,
			Name: uuid.NewString(),
		}

		testRepository := NewUserRepo(testSqlDB)
		err := testRepository.Create(user)
		assert.Nil(t, err)

		gotUser, err := testRepository.Get(userId)
		assert.Nil(t, err)

		assert.Equal(t, &user, gotUser)
	})
}
