package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var DB *sql.DB

func SetupDB() func() {
	var err error
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("AWS_MYSQL_USERNAME"),
		os.Getenv("AWS_MYSQL_PASSWORD"),
		os.Getenv("AWS_MYSQL_ENDPOINT"),
		os.Getenv("AWS_MYSQL_DBNAME"),
	)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	return func() {
		err = DB.Close()
		if err != nil {
			log.Println(err)
		} else {
			log.Println("db closed")
		}
	}
}
