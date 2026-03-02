package db

import (
	"log"
	"database/sql"
)

import _ "github.com/ncruces/go-sqlite3/driver"
import _ "github.com/ncruces/go-sqlite3/embed"

var sqlDB *sql.DB

func InitDB(){
	sqlDB, err := sql.Open("sqlite3", "file:app.db")
	if(err != nil) {
		log.Fatal(err)
	} 
	_, err = sqlDB.Exec(`CREATE TABLE users (id INT, name VARCHAR(10))`)
	if err != nil {
		log.Fatal(err)
	}
}

