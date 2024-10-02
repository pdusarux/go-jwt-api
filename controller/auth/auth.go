package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pdusarux/go-jwt-api/orm"
	"golang.org/x/crypto/bcrypt"
)

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
