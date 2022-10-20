package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const driver = "mysql"
const connMaxLifeTimeSeconds = 0 //0 means infinite
const maxIdleConnections = 50
const maxOpenConnections = 50

type Connection struct {
	db *sql.DB
}

func NewConnection(dbHost string, dbPort string, dbUser string, dbPass string, dbName string) *Connection {
	db, openError := sql.Open(driver, buildDsn(dbHost, dbPort, dbUser, dbPass, dbName))
	if openError != nil {
		panic(fmt.Sprintf("unable to open db connection at: '%s:%s/%s' due to: %s", dbHost, dbPort, dbName, openError))
	}

	db.SetConnMaxLifetime(connMaxLifeTimeSeconds * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)

	return &Connection{
		db: db,
	}
}

func (c *Connection) NewQueryBuilder() *QueryBuilder {
	return NewQueryBuilder(c.db)
}

func (c *Connection) Execute(query string, params ...interface{}) (sql.Result, error) {
	result, err := c.db.Exec(query, params...)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL Execute method - query: %s", err, query))
	}

	return result, err
}

func (c *Connection) Query(query string, params ...interface{}) (*sql.Rows, error) {
	rows, err := c.db.Query(query, params...)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL Query method - query: %s", err, query))
	}

	return rows, err
}

func (c *Connection) StartTransaction() (*sql.Tx, error) {
	tx, err := c.db.Begin()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL StartTransaction method", err))
	}

	return tx, err
}

func (c *Connection) CommitTransaction(transaction *sql.Tx) error {
	err := transaction.Commit()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL CommitTransaction method", err))
	}

	return err
}

func (c *Connection) RollbackTransaction(transaction *sql.Tx) error {
	err := transaction.Rollback()
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL RollbackTransaction method", err))
	}

	return err
}

func (c *Connection) ExecuteWithTransaction(tx *sql.Tx, query string, params ...interface{}) (sql.Result, error) {
	result, err := tx.Exec(query, params...)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL ExecuteWithTransaction method - query: %s", err, query))
	}

	return result, err
}

func (c *Connection) QueryWithTransaction(tx *sql.Tx, query string, params ...interface{}) (*sql.Rows, error) {
	result, err := tx.Query(query, params...)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error %s when running SQL QueryWithTransaction method - query: %s", err, query))
	}

	return result, err
}

func buildDsn(dbHost string, dbPort string, dbUser string, dbPass string, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
}
