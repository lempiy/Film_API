package models

import (
	"database/sql"
	"fmt"
	"github.com/lempiy/echo_api/types/db"
	"os"
	_ "github.com/lib/pq"
)

var Database db.Database
var err error

func init() {
	info := os.Getenv("DATABASE_URL")
	if info == "" {
		info = fmt.Sprintf("user=%s password=%s dbname=%s host=postgres port=5432 sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DB"))
	}
	Database.DB, err = sql.Open("postgres", info)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}
