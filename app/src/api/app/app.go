package app

import (
	"api/app/items"
  "api/app/files"
	"database/sql"
	"time"
    "fmt"

	"github.com/gin-gonic/gin"
	// Use mysql driver instead
	_ "github.com/go-sql-driver/mysql"
)

var (
	r *gin.Engine
)

const (
	port string = ":8081"
)

// StartApp ...
func StartApp() {
	r = gin.Default()
	db := configDataBase()
	items.Configure(r, db)
	files.Configure(r, db)
	r.Run(port)
}

func configDataBase() *sql.DB {
	//db, err := sql.Open("sqlite3", "./foo.db")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", "user", "userpwd", "db", "db"))
	if err != nil {
		panic("Could not connect to the db")
	}

	for {
		err := db.Ping()
		if err != nil {
			time.Sleep(1*time.Second)
			continue
		}
		// This is bad practice... You should create a schema.sql with all the definitions
		createItemTable(db)
    createFileTable(db)
		return db
	}

}

func createItemTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS items(
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
		name TEXT,
		description TEXT
	);`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func createFileTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE IF NOT EXISTS files(
		id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
		titulo TEXT,
		descripcion TEXT
	);`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}
