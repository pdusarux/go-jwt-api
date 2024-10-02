package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/pdusarux/go-jwt-api/controller/auth"
	"github.com/pdusarux/go-jwt-api/orm"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	orm.InitDB()

	r.POST("/register", auth.Register)

	r.Run(":8000")
}
