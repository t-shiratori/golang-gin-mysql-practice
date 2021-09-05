package controllers

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type TestController struct{}

type User struct {
	ID    string
	Name  string
	Email string
}

func (c TestController) Get(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var users []User

		rows, err := db.Query("select * from users;")

		defer rows.Close()

		if err != nil {
			fmt.Println(">>> db Query ERROR!")
		}

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(">>> row Scan ERROR!")
			}
			users = append(users, user)
		}

		fmt.Println("users")
		fmt.Printf("%+v\n", users)

		ctx.HTML(200, "index.html", gin.H{
			"users": users,
		})
	}
}

func (c TestController) Add(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		name := ctx.PostForm("name")
		email := ctx.PostForm("email")
		fmt.Println("create user " + name + " with email " + email)

		prepareDB, err := db.Prepare("insert into users (name, email) values(?,?);")

		defer prepareDB.Close()

		if err != nil {
			fmt.Println(">>> db Prepare Insert ERROR!")
		}

		result, err := prepareDB.Exec(name, email)

		if err != nil {
			fmt.Println(">>> db Insert ERROR!")
		}

		fmt.Println(result)

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			fmt.Println(">>> db get lastInsertID ERROR!")
		}
		fmt.Println(lastInsertID)

		ctx.Redirect(302, "/test")
	}
}

func (c TestController) Delete(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		param := ctx.Param("id")
		id, err := strconv.Atoi(param)

		prepareDB, err := db.Prepare("delete from users where id=?;")

		if err != nil {
			fmt.Println(">>> db Prepare delete ERROR!")
		}

		defer prepareDB.Close()

		result, err := prepareDB.Exec(id)

		if err != nil {
			fmt.Println(">>> db delete ERROR!")
		}

		fmt.Println(result)

		ctx.Redirect(302, "/test")
	}
}
