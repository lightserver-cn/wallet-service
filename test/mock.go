package test

import (
	"database/sql"
	"log"
	"net/http/httptest"
	"testing"

	"server/router"
	"server/test/db"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MockTest struct {
	Expect    *httpexpect.Expect
	CleanFunc []func() error
}

func NewMockTest() *MockTest {
	return &MockTest{}
}

func getExpect(t *testing.T, sqlDB *sql.DB, logger *zap.SugaredLogger) *httpexpect.Expect {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router.Router(engine, sqlDB, logger)
	server := httptest.NewServer(engine)
	return httpexpect.New(t, server.URL)
}

type postgresqlConf struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

func (m *MockTest) start(t *testing.T) *MockTest {
	dbConf := postgresqlConf{
		Driver:   "postgres",
		Host:     "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "test_postgres",
	}

	dbTest, err := db.NewDalTest(dbConf.Driver, dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
	if err != nil {
		log.Fatalf("db.NewDalTest err: %v", err)
	}

	err = dbTest.DropTestDB()
	if err != nil {
		log.Fatalf("DropTestDB err: %v", err)
	}

	err = dbTest.CreateTestDB()
	if err != nil {
		log.Fatalf("CreateTestDB err: %v", err)
	}

	err = dbTest.CreateTestTables()
	if err != nil {
		log.Fatalf("CreateTestTables err: %v", err)
	}

	err = dbTest.CreateTestData()
	if err != nil {
		log.Fatalf("CreateTestData err: %v", err)
	}

	m.CleanFunc = append(m.CleanFunc, dbTest.TruncateTable, dbTest.DropTestDB, dbTest.Close)

	m.Expect = getExpect(t, dbTest.DB(), zap.NewExample().Sugar())

	return m
}

func (m *MockTest) Teardown() {
	m.Expect = nil
	for _, f := range m.CleanFunc {
		_ = f()
	}
}
