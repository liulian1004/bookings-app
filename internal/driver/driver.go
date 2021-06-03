package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4"
)


type DB struct {
	SQL *sql.DB
}

var dbCon = &DB{}


const maxOpenDbCon = 10 //max 10 connects for db
const maxIDleDbCon = 5
const maxDbLifetime = 5 * time.Minute // max time to connect to db

func ConnectDB(dsn string) (*DB, error) {
	// create new db and set up config
	d, err := NewDataBase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbCon)
	d.SetConnMaxLifetime(maxDbLifetime)
	d.SetMaxIdleConns(maxIDleDbCon)

	dbCon.SQL = d

	err =testDB(d) // ping db
	
	if err != nil {
		return nil, err
	}
	return dbCon, nil

}

func testDB(d *sql.DB) error {
	err := d.Ping();
	if err != nil {
		return err
	}
	return nil
}
func NewDataBase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}