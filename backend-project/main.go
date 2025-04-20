package main

import (
	"backend-project/config"
	"backend-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.ConnectDatabase()

	routes.SetupRoutes(r)

	r.Static("/uploads", "./uploads")

	r.Run(":8080")
}
