package db

import (
	"testing"

	_ "github.com/lib/pq"
	"go.uber.org/goleak"
)

const (
	testDriver        = "postgres"
	testHost          = "localhost"
	testPort          = 5432
	testUser          = "postgres"
	testPassword      = "postgres"
	testDBName        = "db_test"
	testDefaultDbname = "test_postgres" // when drop database
	testTempDbname    = "template0"
)

func TestNewDalTest(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	if dal == nil {
		t.Errorf("NewDalTest returned nil")
	}
}

func TestNewDalTest_ErrNotCurrentlyAcceptingConnections(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testTempDbname)
	if err == nil {
		t.Errorf("NewDalTest did not return an error when not currently accepting connections")
	}

	if dal != nil {
		_ = dal.Close()
		t.Errorf("NewDalTest returned a non-nil value when not currently accepting connections")
	}
}

func TestNewDalTest_ConnectionFailure(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest("postgres", "wrong_host", testPort, testUser, testPassword, testDBName)
	if err == nil {
		t.Errorf("NewDalTest did not return an error when connection failure")
	}

	if dal != nil {
		t.Errorf("NewDalTest returned a non-nil value when connection should have failed")
	} else {
		t.Logf("NewDalTest correctly returned nil when connection failed")
	}

	expectedError := "dial tcp: lookup wrong_host: no such host"
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expected error: %q, got: %q", expectedError, err.Error())
	}
}

func TestClose(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}

	err = dal.Close()
	if err != nil {
		t.Errorf("Close returned an error: %v", err)
	}
}

func TestClose_ErrClose(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}

	err = dal.Close()
	if err != nil {
		t.Errorf("Close returned an error: %v", err)
	}
}

func TestCreateTestDB(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	// CreateTestDB: failed to create database: pq: duplicate key value violates unique constraint "pg_database_datname_index
	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	rows, err := dal.DB().Query("SELECT datname FROM pg_database WHERE datname = 'db_test'")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("Database was not created")
	}
}

func TestCreateTestDB_AlreadyExists(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer func(dal *DalTest) {
		err = dal.Close()
	}(dal)

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}
}

func TestDropTestDB(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	rows, err := dal.DB().Query("SELECT datname FROM pg_database WHERE datname = 'db_test'")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		t.Errorf("Database was not dropped")
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}
}

func TestDropTestDB_ErrCurrentlyOpenDatabase(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDBName)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err == nil {
		t.Errorf("DropTestDB did not return an error when drop the currently open database")
	}
}

func TestCreateTestTables(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	err = dal.CreateTestTables()
	if err != nil {
		t.Errorf("CreateTestTables returned an error: %v", err)
	}

	rows, err := dal.DB().Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'" +
		" AND table_name IN ('t_user', 't_wallet', 't_transaction')")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	expectedTables := []string{"t_user", "t_wallet", "t_transaction"}
	tablesFound := make(map[string]bool)

	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			t.Errorf("Scan returned an error: %v", err)
		}
		tablesFound[tableName] = true
	}

	for _, table := range expectedTables {
		if !tablesFound[table] {
			t.Errorf("Table %s was not created", table)
		}
	}
}

func TestCreateTestData(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	err = dal.CreateTestTables()
	if err != nil {
		t.Errorf("CreateTestTables returned an error: %v", err)
	}

	err = dal.CreateTestData()
	if err != nil {
		t.Errorf("CreateTestData returned an error: %v", err)
	}

	rows, err := dal.DB().Query("SELECT COUNT(*) FROM t_user")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			t.Errorf("Scan returned an error: %v", err)
		}
	} else {
		t.Errorf("No data found in t_user")
	}

	if count == 0 {
		t.Errorf("No data inserted into t_user")
	}
}

func TestTruncateTable(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"))

	dal, err := NewDalTest(testDriver, testHost, testPort, testUser, testPassword, testDefaultDbname)
	if err != nil {
		t.Errorf("NewDalTest returned an error: %v", err)
	}
	defer dal.Close()

	err = dal.DropTestDB()
	if err != nil {
		t.Errorf("DropTestDB returned an error: %v", err)
	}

	err = dal.CreateTestDB()
	if err != nil {
		t.Errorf("CreateTestDB returned an error: %v", err)
	}

	err = dal.CreateTestTables()
	if err != nil {
		t.Errorf("CreateTestTables returned an error: %v", err)
	}

	err = dal.CreateTestData()
	if err != nil {
		t.Errorf("CreateTestData returned an error: %v", err)
	}

	rows, err := dal.DB().Query("SELECT COUNT(*) FROM t_user")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	var countBefore int
	if rows.Next() {
		if err = rows.Scan(&countBefore); err != nil {
			t.Errorf("Scan returned an error: %v", err)
		}
	} else {
		t.Errorf("No data found in t_user")
	}

	if countBefore == 0 {
		t.Errorf("No data inserted into t_user before truncation")
	}

	err = dal.TruncateTable()
	if err != nil {
		t.Errorf("TruncateTable returned an error: %v", err)
	}

	rows, err = dal.DB().Query("SELECT COUNT(*) FROM t_user")
	if err != nil {
		t.Errorf("Query returned an error: %v", err)
	}
	defer rows.Close()

	var countAfter int
	if rows.Next() {
		if err = rows.Scan(&countAfter); err != nil {
			t.Errorf("Scan returned an error: %v", err)
		}
	} else {
		t.Errorf("No data found in t_user after truncation")
	}

	if countAfter != 0 {
		t.Errorf("Data was not truncated from t_user")
	}
}
