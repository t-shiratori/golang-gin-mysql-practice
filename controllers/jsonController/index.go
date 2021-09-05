package jsonController

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type JsonController struct{}

type User struct {
	ID    string
	Name  string
	Email string
}

type JsonAddRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type JsonUpdateRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type JsonDeleteRequest struct {
	ID string `json:"id"`
}

func (c JsonController) Get(db *sql.DB) func(ctx *gin.Context) {
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

func (c JsonController) Add(db *sql.DB) func(ctx *gin.Context) {
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

func (c JsonController) Update(db *sql.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		var jsonReq JsonUpdateRequest

		if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("update user name " + jsonReq.Name + " with email " + jsonReq.Email)

		prepareDB, err := db.Prepare("update users set name = ?, email = ? where id = ?;")

		defer prepareDB.Close()

		if err != nil {
			fmt.Println(">>> db Prepare Update ERROR!")
		}

		convertedId, err := strconv.ParseInt(jsonReq.ID, 10, 64)

		result, err := prepareDB.Exec(jsonReq.Name, jsonReq.Email, convertedId)

		if err != nil {
			fmt.Println(">>> db Update ERROR!")
		}

		fmt.Println(result)

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fmt.Println(">>> db get RowsAffected ERROR!")
		}
		fmt.Println("rowsAffected: " + strconv.FormatInt(rowsAffected, 10))

		ctx.JSON(http.StatusOK, gin.H{
			"user": jsonReq,
		})
	}
}

func (c JsonController) Delete(db *sql.DB) func(ctx *gin.Context) {
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
