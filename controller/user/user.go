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
