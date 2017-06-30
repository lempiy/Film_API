package main

import (
	"database/sql"
	"fmt"
	"github.com/lempiy/echo_api/types/db"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"strings"
)

var database db.Database
var err error

func init() {
	info := os.Getenv("DATABASE_URL")
	if info == "" {
		info = fmt.Sprintf("user=%s password=%s dbname=%s host=postgres port=5432 sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DB"))
	}
	database.DB, err = sql.Open("postgres", info)
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <path_to_file>.sql\n", os.Args[0])
		return
	}
	path := os.Args[1]
	if !strings.HasSuffix(path, ".sql") {
		fmt.Println("Only *.sql files allowed.")
	}
	fileBuffer, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error during reading file - %s\n%v", path, err)
		return
	}
	sqlScript := string(fileBuffer)
	_, err = database.Exec(sqlScript)
	if err != nil {
		fmt.Printf("Error while executing SQL - %s\n%v", path, err)
		return
	}
	fmt.Println("Initailization successful.")
	return
}
