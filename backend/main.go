package main

import (
	"ecommerce/database"
	"ecommerce/helpers"
	"ecommerce/middlewares"
	"ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func init() {

	dbHandler := &database.Mongo{}

	dbHandler.ConnectDB()

	helpers.DBHelper = dbHandler
	routes.DBHelper = dbHandler
	middlewares.DBHelper = dbHandler

}

func main() {

	router := gin.Default()

	router.POST("/api/v1/auth/register", routes.Register)
	router.POST("/api/v1/auth/login", routes.Login)
	router.GET("/test", middlewares.Authenticate, routes.Test)
	router.GET("/admin", middlewares.Authenticate, middlewares.IsAdmin, routes.Test)

	router.Run()

}
