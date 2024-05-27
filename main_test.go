package frog_blossom_db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/frog-blossom-db?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	//conn, err := sql.Open(dbDriver, dbSource)
	//if err != nil {
	//	log.Fatal("cannot connect to database:", err)
	//}

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
