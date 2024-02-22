package basicgorm

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connection(database string) (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error al obtener configuración del servidor")
	}
	server := os.Getenv("ENV_DDBB_SERVER")
	user := os.Getenv("ENV_DDBB_USER")
	password := os.Getenv("ENV_DDBB_PASSWORD")
	port := os.Getenv("ENV_DDBB_PORT")
	ssl := os.Getenv("ENV_DDBB_SSL")
	sslmode := "disable"
	if ssl == "true" {
		sslmode = "require"
	}
	connection_string := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", server, port, user, password, database, sslmode)
	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		errs := fmt.Errorf("error connection: %s ", err.Error())
		return nil, errs
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		errs := fmt.Errorf("error creating connection: %s", err.Error())
		return nil, errs
	}
	return db, nil
}

func ConnectionCloud() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("error al obtener configuración del servidor cloud")
	}
	server := os.Getenv("ENV_DDBB_SERVER")
	user := os.Getenv("ENV_DDBB_USER")
	password := os.Getenv("ENV_DDBB_PASSWORD")
	database := os.Getenv("ENV_DDBB_DATABASE")
	port := os.Getenv("ENV_DDBB_PORT")
	ssl := os.Getenv("ENV_DDBB_SSL")
	sslmode := "disable"
	if ssl == "true" {
		sslmode = "require"
	}

	connection_string := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", server, port, user, password, database, sslmode)

	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		errs := fmt.Errorf("error connection: %s ", err.Error())
		return nil, errs
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		errs := fmt.Errorf("error creating connection: %s", err.Error())
		return nil, errs
	}
	return db, nil
}
