package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries

const driverName = "postgres"
const dataSourceName = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

func TestMain(m *testing.M) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("failed TESTMAIN : %+v", err)
	}

	TestQueries = New(db)

	os.Exit(m.Run())

}
