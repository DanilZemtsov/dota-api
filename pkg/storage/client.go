package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// root:@tcp(127.0.0.1:3306)/golang
func Connect(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
