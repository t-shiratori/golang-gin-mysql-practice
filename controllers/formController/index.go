package formController

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type FormController struct{}

type User struct {
	ID    string
	Name  string
	Email string
}

var rootPash = "/form"
var editPash = "/form/edit"

func (c FormController) Get(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var users []User

		rows, err := db.Query("select * from users;")

		defer rows.Close()

		if err != nil {
			fmt.Println("--> db.Query() error")
		}

		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println("--> row.Scan() error")
			}
			users = append(users, user)
		}

		fmt.Printf("users %+v\n", users)

		ctx.HTML(200, "index.html", gin.H{
			"users": users,
		})
	}
}

func (c FormController) Add(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var user User
		user.Name = ctx.PostForm("name")
		user.Email = ctx.PostForm("email")

		fmt.Printf("user %+v\n", user)

		prepareDB, err := db.Prepare("insert into users (name, email) values(?,?);")

		defer prepareDB.Close()

		if err != nil {
			fmt.Println("--> db.Prepare() error")
		}

		result, err := prepareDB.Exec(user.Name, user.Email)

		if err != nil {
			fmt.Println("--> db insert error")
		}

		fmt.Println(result)

		lastInsertID, err := result.LastInsertId()

		if err != nil {
			fmt.Println("--> result.LastInsertId() error")
		}

		fmt.Println(lastInsertID)

		ctx.Redirect(302, rootPash)
	}
}

func (c FormController) Delete(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var user User
		user.ID = ctx.Param("id")

		prepareDB, err := db.Prepare("delete from users where id=?;")

		if err != nil {
			fmt.Println("--> db.Prepare() error")
		}

		defer prepareDB.Close()

		result, err := prepareDB.Exec(user.ID)

		if err != nil {
			fmt.Println("--> db delete error")
		}

		fmt.Println(result)

		ctx.Redirect(302, rootPash)
	}
}

func (c FormController) Edit(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var user User
		user.ID = ctx.Param("id")

		row := db.QueryRow("select * from users where id = ?;", user.ID)

		err := row.Scan(&user.ID, &user.Name, &user.Email)

		if err != nil {
			fmt.Println("--> db Scan error")
		}

		ctx.HTML(200, "edit.html", gin.H{
			"user": user,
		})
	}
}

func (c FormController) Update(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var user User
		user.ID = ctx.Param("id")
		user.Name = ctx.PostForm("name")
		user.Email = ctx.PostForm("email")

		fmt.Printf("user %+v\n", user)

		prepareDB, err := db.Prepare("update users set name = ?, email = ? where id = ?;")

		if err != nil {
			fmt.Println("--> db.Prepare() error")
		}

		defer prepareDB.Close()

		result, err := prepareDB.Exec(user.Name, user.Email, user.ID)

		if err != nil {
			fmt.Println("--> db update error")
		}

		fmt.Println(result)

		ctx.Redirect(302, rootPash)
	}
}
