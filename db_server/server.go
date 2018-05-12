package main

import (
	"database/sql"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type user struct {
	ID       string `json:"id"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Fullname string `json:"full_name" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}

func main() {
	db := connectDB()
	runServer(db)
}

func connectDB() *sql.DB {
	connStr := "user=postgres password=postgres dbname=instagram_development sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func runServer(db *sql.DB) {
	router := gin.Default()
	router.GET("/user/*id", handleGetReq(db))
	router.POST("/user", handleSetUser(db))
	router.Run()
}

func handleGetReq(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var id = strings.Replace(c.Param("id"), "/", "", 1)

		user, err := getUser(db)(id)

		if err != nil {
			c.JSON(200, gin.H{
				"error": "Internal Server Error",
			})
			return
		}

		c.JSON(200, convertUserToJSON(user))
	}
}

func handleSetUser(db *sql.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		var newUser user

		c.BindJSON(&newUser)

		savedUser, err := setUser(db)(newUser)

		if err != nil {
			c.JSON(200, gin.H{
				"error": err.Detail,
			})
			return
		}

		c.JSON(200, convertUserToJSON(savedUser))
	}
}

func getUser(db *sql.DB) func(string) (user, error) {
	return func(id string) (user, error) {
		var user user

		rows, _ := db.Query("SELECT id, username, email, gender, full_name  FROM users WHERE id = $1", id)

		rows.Next()
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Gender, &user.Fullname)

		return user, err
	}
}

func setUser(db *sql.DB) func(user) (user, *pq.Error) {
	return func(newUser user) (user, *pq.Error) {
		var savedUser user

		id := GenerateUUID()
		err := db.
			QueryRow(`INSERT INTO 
			users(id, username, email, password_hash, is_private, full_name, gender, created_at) 
			VALUES($1,$2, $3, 'hey', FALSE, $4, $5, '2018-03-13 19:15:27.512')
			RETURNING id, username, email, full_name, gender
			`, id, newUser.Username, newUser.Email, newUser.Fullname, newUser.Gender).
			Scan(&savedUser.ID, &savedUser.Username, &savedUser.Email, &savedUser.Fullname, &savedUser.Gender)
		if err != nil {
			return savedUser, err.(*pq.Error)
		}

		return savedUser, nil
	}
}

func convertUserToJSON(user user) gin.H {
	return gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.Fullname,
			"gender":    user.Gender,
		},
	}
}

// GenerateUUID Generates a new UUID
func GenerateUUID() string {
	uuid, _ := uuid.NewUUID()
	return uuid.String()
}
