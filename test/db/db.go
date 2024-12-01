package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func NewDalTest(driver, host string, port int64, user, password, dbName string) (*DalTest, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		log.Printf("NewDalTest: failed to open database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("NewDalTest: failed to ping database: %v", err)
		return nil, err
	}

	return &DalTest{db: db}, nil
}

type DalTest struct {
	db *sql.DB
}

func (d *DalTest) Close() error {
	return d.db.Close()
}

func (d *DalTest) DB() *sql.DB {
	return d.db
}

func (d *DalTest) DBExists(dbName string) (bool, error) {
	var count int

	err := d.db.QueryRow(fmt.Sprintf(`SELECT 1 FROM pg_catalog.pg_database WHERE datname = '%s'`, dbName)).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to query for database existence: %w", err)
	}

	return count > 0, nil
}

func (d *DalTest) CreateTestDB() error {
	var err error
	exists, err := d.DBExists("db_test")
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if exists {
		if err = d.DropTestDB(); err != nil {
			return fmt.Errorf("failed to drop existing database: %w", err)
		}
		log.Println("CreateTestDB: existing database dropped successfully")
	}

	_, err = d.db.Exec(`CREATE DATABASE db_test`)
	if err != nil {
		log.Printf("CreateTestDB: failed to create database: %v", err)
		return err
	}

	log.Println("CreateTestDB: database created successfully")
	return nil
}

func (d *DalTest) DropTestDB() error {
	_, err := d.db.Exec(`DROP DATABASE IF EXISTS db_test`)
	if err != nil {
		log.Printf("DropTestDB: failed to drop database: %v", err)
		return err
	}

	log.Println("DropTestDB: database dropped successfully")
	return nil
}

func (d *DalTest) CreateTestTables() error {
	dir, err := GetDirPath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(filepath.Join(dir, ".", "ddl.sql"))
	if err != nil {
		log.Printf("CreateTestTables: failed to read SQL file: %v", err)
		return err
	}

	_, err = d.db.Exec(string(content))
	if err != nil {
		log.Printf("CreateTestTables: failed to execute SQL: %v", err)
		return err
	}

	log.Println("CreateTestTables: tables created successfully")
	return nil
}

func (d *DalTest) CreateTestData() error {
	dir, errPath := GetDirPath()
	if errPath != nil {
		return errPath
	}

	sqlFiles := []string{
		"user.sql",
		"wallet.sql",
		"transaction.sql",
	}

	var combinedContent string
	for _, sqlFile := range sqlFiles {
		content, err := os.ReadFile(filepath.Join(dir, ".", sqlFile))
		if err != nil {
			log.Printf("Failed to read SQL file %s: %v", sqlFile, err)
			return err
		}

		combinedContent += string(content)
	}

	_, err := d.db.Exec(combinedContent)
	if err != nil {
		log.Printf("Failed to execute SQL : %v", err)
		return err
	}

	log.Println("CreateTestData: Successfully executed SQL")
	return nil
}

func (d *DalTest) TruncateTable() error {
	tables := []string{
		"t_user",
		"t_wallet",
		"t_transaction",
	}

	tx, err := d.db.Begin()
	if err != nil {
		log.Printf("TruncateTable: failed to begin transaction: %v", err)
		return err
	}

	for _, table := range tables {
		_, err = tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			_ = tx.Rollback()
			log.Printf("TruncateTable: failed to truncate table %s: %v", table, err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("TruncateTable: failed to commit transaction: %v", err)
		return err
	}

	log.Println("TruncateTable: tables truncated successfully")
	return nil
}
