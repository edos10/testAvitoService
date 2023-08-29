package databases

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func CreateDatabaseConnect() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	fmt.Println(dbHost, dbPort, dbUser, dbPassword, dbName, "quuuu")
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	return db, err
}
