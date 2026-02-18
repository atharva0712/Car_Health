package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dsn := "root:password@tcp(127.0.0.1:3306)/car_monitor"  // data source name (connection string)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Databse Connected")
	return db

}
