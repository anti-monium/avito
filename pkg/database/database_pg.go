package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func ConnectDatabase() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	pass := os.Getenv("POSTGRES_PASSWORD")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	log.Info(psqlSetup)
	db, errSql := sql.Open("postgres", psqlSetup)
	if errSql != nil {
		log.Info(fmt.Sprint(errSql.Error()))
		return nil, errSql
	}
	log.Info(fmt.Sprint(db))
	return db, nil
}
