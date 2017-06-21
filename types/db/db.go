package db

import (
	"database/sql"
	"fmt"
)

type Database struct {
	DB *sql.DB
}

func (db Database) begin() (transaction *sql.Tx) {
	transaction, err := db.DB.Begin()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return transaction
}

func (db Database) prepare(query string) (statement *sql.Stmt) {
	statement, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return statement
}

func (db Database) Exec(query string, arguments ...interface{}) (result sql.Result, err error) {
	result, err = db.DB.Exec(query, arguments...)
	return
}

func (db Database) Query(query string, arguments ...interface{}) (rows *sql.Rows) {
	rows, err := db.DB.Query(query, arguments...)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}

func (db Database) queryRow(query string, arguments ...interface{}) (row *sql.Row) {
	row = db.DB.QueryRow(query, arguments...)
	return row
}

//singleQuery multiple query isolation
func (db Database) SingleQuery(sql string, args ...interface{}) error {
	SQL := db.prepare(sql)
	tx := db.begin()
	_, err := tx.Stmt(SQL).Exec(args...)

	if err != nil {
		fmt.Println("singleQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("singleQuery successful")
	}
	return err
}

//InsertWithReturningID exec statement and return inserted ID
func (db Database) InsertWithReturningID(sql string, args ...interface{}) int {
	var lastID int64
	sql = sql + " RETURNING id;"

	row := db.queryRow(sql, args...)
	row.Scan(&lastID)

	id := int(lastID)
	fmt.Printf("insertWithReturningID: %d\n", id)
	return id
}

//singleQueryWithAffected multiple query isolation returns affected rows
func (db Database) SingleQueryWithAffected(sql string, args ...interface{}) (int, error) {
	SQL := db.prepare(sql)
	tx := db.begin()
	result, err := tx.Stmt(SQL).Exec(args...)

	affectedCount, err := result.RowsAffected()
	id := int(affectedCount)
	if err != nil {
		fmt.Println("singleQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		fmt.Println("singleQuery successful")
	}
	return id, err
}

//Close func closes DB connection
func (db Database) Close() {
	db.DB.Close()
}
