package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	db   *sql.DB
	once sync.Once
)

func Connect() *sql.DB {
	// dsn := "root:password@tcp(127.0.0.1:3306)/car_monitor?parseTime=true" // data source name (connection string)

	// db, err := sql.Open("mysql", dsn)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = db.Ping()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Databse Connected")
	// return db
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
			log.Fatal("Database environment variables are not properly set")
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

		database, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("Failed to Open DB", err)
		}

		err1 := database.Ping()
		if err1 != nil {
			log.Fatal("Failed to connect to DB", err1)
		}

		log.Println("Database connected")

		db = database
	})
	return db
}
