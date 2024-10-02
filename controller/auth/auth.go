package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pdusarux/go-jwt-api/orm"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	// Check Username Exists
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User Exists"})
		return
	}

	// Encrypt Password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Declare variables and insert data
	user := orm.User{Username: json.Username, Password: string(encryptedPassword),
		Fullname: json.Fullname, Avatar: json.Avatar}

	// Save in db
	if result := orm.Db.Create(&user); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Return result
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Registered", "userId": user.ID})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User Create Failed"})
	}

}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First(&userExist)
	if userExist.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User Does Not Exists"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Login Failed"})
		return
	}

	hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userExist.ID,
		"exp":    time.Now().Add(time.Minute * 1).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Token is wrong"})
		return
	}

	fmt.Println(tokenString, err)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Login Success", "token": tokenString})

}
