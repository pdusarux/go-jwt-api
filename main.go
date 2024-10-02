package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type User struct {
	gorm.Model
	Username string
	Password string
	Fullname string
	Avatar   string
}

func main() {
	dsn := "root:root_password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", func(c *gin.Context) {
		var json Register
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		// Check Username Exists
		var userExist User
		db.Where("username = ?", json.Username).First(&userExist)
		if userExist.ID > 0 {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Exists"})
			return
		}

		// Encrypt Password
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Declare variables and insert data
		user := User{Username: json.Username, Password: string(encryptedPassword),
			Fullname: json.Fullname, Avatar: json.Avatar}

		// Save in db
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// Return result
		if user.ID > 0 {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Registered", "userId": user.ID})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "error", "message": "User Create Failed"})
		}

	})

	r.Run(":8000")
}
