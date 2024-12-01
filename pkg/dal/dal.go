package dal

import (
	"database/sql"
	"log"

	"github.com/redis/go-redis/v9"
)

var CustomDal *Dal

type Dal struct {
	DB     *sql.DB
	RDB    redis.UniversalClient
	logger *log.Logger
}

func New(db *sql.DB, rdb *redis.Client, logger *log.Logger) (*Dal, error) {
	return &Dal{DB: db, RDB: rdb, logger: logger}, nil
}
