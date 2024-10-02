package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pdusarux/go-jwt-api/orm"
)

func ReadAll(c *gin.Context) {
	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User Read Success",
		"users":   users,
	})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user orm.User
	orm.Db.Find(&user, userId)
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Profile Read Success",
		"profile": user,
	})
}
