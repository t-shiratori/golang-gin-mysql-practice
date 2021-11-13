package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	user := "go_test"
	password := "password"
	ip := "db" // docker-composeのサービス名
	port := "3306"
	dbName := "go_database"
	option := "?charset=utf8mb4&parseTime=True&loc=Local"

	// 詳細は https://github.com/go-sql-driver/mysql#dsn-data-source-name を参照
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName + option
	db, err := sql.Open("mysql", dsn)

	fmt.Println("--> db open")

	if err != nil {
		fmt.Println("--> db open error")
	}

	fmt.Println("--> db send ping")

	err = db.Ping()

	if err != nil {
		fmt.Println("--> db ping error")
	}

	return db
}
