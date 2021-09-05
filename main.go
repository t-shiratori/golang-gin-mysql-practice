package main

import (
	"fmt"
	"gin-api/controllers/formController"
	"gin-api/controllers/jsonController"
	"gin-api/driver"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	formController := formController.FormController{}
	jsonController := jsonController.JsonController{}

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

	// form endpoints
	router.GET("/test/", formController.Get(db))
	router.POST("/test/add", formController.Add(db))
	router.POST("/test/delete/:id", formController.Delete(db))

	// json endpoints
	router.GET("/user/", jsonController.Get(db))
	router.POST("/user/", jsonController.Add(db))
	router.PUT("/user/", jsonController.Update(db))
	router.DELETE("/user/", jsonController.Delete(db))

	router.Run()
}
