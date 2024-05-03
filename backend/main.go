package main

import (
	"ecommerce/database"
	"ecommerce/helpers"
	"ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func init() {

	dbHandler := &database.Mongo{}

	dbHandler.ConnectDB()

	helpers.DBHelper = dbHandler
	routes.DBHelper = dbHandler

}

func main() {

	router := gin.Default()

	router.POST("/api/v1/auth/register", routes.Register)

	router.Run()

}
