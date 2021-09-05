package userController

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type UserController struct{}

type User struct {
	ID    string
	Name  string
	Email string
}

type JsonAddRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type JsonDeleteRequest struct {
	ID string `json:"id"`
}

func (c UserController) Get(db *sql.DB) func(ctx *gin.Context) {
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

		ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}

func (c UserController) Add(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var user User

		var jsonReq JsonAddRequest

		if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("create user " + jsonReq.Name + " with email " + jsonReq.Email)

		prepareDB, err := db.Prepare("insert into users (name, email) values(?,?);")

		defer prepareDB.Close()

		if err != nil {
			fmt.Println(">>> db Prepare Insert ERROR!")
		}

		result, err := prepareDB.Exec(jsonReq.Name, jsonReq.Email)

		if err != nil {
			fmt.Println(">>> db Insert ERROR!")
		}

		fmt.Println(result)

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			fmt.Println(">>> db get lastInsertID ERROR!")
		}
		fmt.Println(lastInsertID)

		user.ID = strconv.FormatInt(lastInsertID, 10)
		user.Name = jsonReq.Name
		user.Email = jsonReq.Email

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (c UserController) Delete(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var jsonReq JsonDeleteRequest

		if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prepareDB, err := db.Prepare("delete from users where id=?;")

		if err != nil {
			fmt.Println(">>> db Prepare delete ERROR!")
		}

		defer prepareDB.Close()

		result, err := prepareDB.Exec(jsonReq.ID)

		if err != nil {
			fmt.Println(">>> db delete ERROR!")
		}

		fmt.Println(result)

		ctx.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	}
}
