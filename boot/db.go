package boot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"server/config"
	"server/pkg/dal"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	ddlPath      = "/usr/local/config/ddl.sql"
	ddlLocalPath = "config/ddl.sql"
)

func initDB() error {
	var err error

	env := os.Getenv("ENV")
	if env == "" {
		ddlPath = ddlLocalPath
	}

	dbConf := config.Config.DB
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)

	db, err := sql.Open(dbConf.Driver, dataSourceName)
	if err != nil {
		return err
	}

	// for test
	_, _ = db.Exec("DROP DATABASE IF EXISTS test_postgres")
	_, _ = db.Exec("CREATE DATABASE test_postgres")
	_, _ = db.Exec("DROP DATABASE IF EXISTS db_test")
	_, _ = db.Exec("CREATE DATABASE db_test")

	if dbConf.InitTable {
		var content []byte
		content, err = os.ReadFile(ddlPath)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}

		log.Printf("------ InitTable Success \n")
	}

	rdbConf := config.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbConf.Addr,
		Password: rdbConf.Password,
		DB:       rdbConf.DB,
		PoolSize: 100,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	logger := log.Default()
	dal.CustomDal, err = dal.New(db, rdb, logger)
	if err != nil {
		return err
	}

	return nil
}
