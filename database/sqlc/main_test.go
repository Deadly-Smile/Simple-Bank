package sqlc

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Initialize the test database connection and queries
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}
	testQueries = New(testDB)

	// Run the tests
	os.Exit(m.Run())
}
