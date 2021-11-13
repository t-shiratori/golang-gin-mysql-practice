package formController

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type FormController struct{}

type User struct {
	ID    string
	Name  string
	Email string
}

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

		ctx.Redirect(302, "/test")
	}
}

func (c FormController) Delete(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		param := ctx.Param("id")
		id, err := strconv.Atoi(param)

		prepareDB, err := db.Prepare("delete from users where id=?;")

		if err != nil {
			fmt.Println("--> db.Prepare() error")
		}

		defer prepareDB.Close()

		result, err := prepareDB.Exec(id)

		if err != nil {
			fmt.Println("--> db delete error")
		}

		fmt.Println(result)

		ctx.Redirect(302, "/test")
	}
}
