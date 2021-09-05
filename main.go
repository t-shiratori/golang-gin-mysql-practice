package main

import (
	"fmt"
	"gin-api/controllers/testController"
	"gin-api/controllers/userController"
	"gin-api/driver"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	testController := testController.TestController{}
	userController := userController.UserController{}

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

	// test endpoints
	router.GET("/test/", testController.Get(db))
	router.POST("/test/add", testController.Add(db))
	router.POST("/test/delete/:id", testController.Delete(db))

	// user endpoints
	router.GET("/user/", userController.Get(db))
	router.POST("/user/", userController.Add(db))
	router.DELETE("/user/", userController.Delete(db))

	router.Run()
}
