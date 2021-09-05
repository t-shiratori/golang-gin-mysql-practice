package main

import (
	"fmt"
	"gin-api/controllers"
	"gin-api/driver"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

//[Go言語でデータベース（MySQL）に接続する方法 | Enjoy IT Life](https://nishinatoshiharu.com/connect-go-database/)

type User struct {
	ID    string
	Name  string
	Email string
}

// type Error struct {
// 	Message string `json:"message"`
// }

// func SendError(w http.ResponseWriter, status int, err Error) {
// 	w.WriteHeader(status)
// 	json.NewEncoder(w).Encode(err)
// }

// func SendSuccess(w http.ResponseWriter, data interface{}) {
// 	json.NewEncoder(w).Encode(data)
// }

func main() {

	testController := controllers.TestController{}

	db := driver.ConnectDB()

	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	fmt.Println(">>> create table!")
	_, err := db.Exec(`
		create table if not exists users (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
			name VARCHAR(255),
			email VARCHAR(255)
		);`)

	if err != nil {
		fmt.Println(">>> create table ERROR!")
	}

	router.GET("/test/", testController.Get(db))

	router.POST("/test/add", testController.Add(db))

	router.POST("/test/delete/:id", testController.Delete(db))

	router.Run()
}
