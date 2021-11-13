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
	router.Static("/assets", "./assets")

	fmt.Println("--> create table")

	_, err := db.Exec(`
		create table if not exists users (
			id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
			name VARCHAR(255),
			email VARCHAR(255)
		);`)

	if err != nil {
		fmt.Println("--> create table error")
	}

	// form endpoints
	router.GET("/form/", formController.Get(db))
	router.GET("/form/edit/:id", formController.Edit(db))
	router.POST("/form/add", formController.Add(db))
	router.POST("/form/delete/:id", formController.Delete(db))
	router.POST("/form/update/:id", formController.Update(db))

	// json endpoints
	router.GET("/api/user/", jsonController.Get(db))
	router.POST("/api/user/", jsonController.Add(db))
	router.PUT("/api/user/", jsonController.Update(db))
	router.DELETE("/api/user/", jsonController.Delete(db))

	router.Run()
}
