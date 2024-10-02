package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/pdusarux/go-jwt-api/controller/auth"
	"github.com/pdusarux/go-jwt-api/controller/user"
	"github.com/pdusarux/go-jwt-api/middleware"
	"github.com/pdusarux/go-jwt-api/orm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	orm.InitDB()

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)

	authorized := r.Group("/users", middleware.JWTAuthen())
	authorized.GET("/readall", user.ReadAll)
	authorized.GET("/profile", user.Profile)

	r.Run(":8000")
}
