package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type DB struct {
	*sql.DB
}

func ReadConfig() string {
	DbDriver := os.Getenv("DB_driver")
	DbUsername := os.Getenv("DB_username")
	DbPassword := os.Getenv("DB_password")
	DbHost := os.Getenv("DB_host")
	DbPort := os.Getenv("DB_port")
	DbName := os.Getenv("DB_name")
	DbSslmode := os.Getenv("DB_sslmode")

	config := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		DbDriver, DbUsername, DbPassword, DbHost, DbPort, DbName, DbSslmode)
	return config
}
func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected!")
	return &DB{db}, nil
}

func Setup(filename string, db *DB) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Setupfile opening error: ", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error after opening setupfile: ", err)
		return
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("Error after opening setupfile: ", err)
		panic(err)
	}

	command := string(bs)
	_, err = db.Exec(command)
	if err != nil {
		fmt.Println("Command error")
	}
}
