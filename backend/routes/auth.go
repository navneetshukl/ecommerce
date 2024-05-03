package routes

import (
	"ecommerce/database"
	"ecommerce/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var DBHelper *database.Mongo

func Register(c *gin.Context) {

	var user models.User

	err := c.BindJSON(&user)
	if err != nil {

		log.Println("Error in reading the body ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in reading the body",
		})
		return

	}

	if len(user.Name) == 0 {
		log.Println("Name is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Name is required",
		})
		return

	}
	if len(user.Email) == 0 {
		log.Println("Email is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is required",
		})
		return

	}
	if len(user.Password) == 0 {
		log.Println("Password is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password is required",
		})
		return

	}
	if len(user.Address) == 0 {
		log.Println("Address is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Address is required",
		})
		return

	}
	if len(user.Phone) == 0 {
		log.Println("Phone is required")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone is required",
		})
		return

	}

	exist, err := DBHelper.CheckUser(user.Email)
	if err != nil {

		log.Println("Error in checking the user ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in checking to database",
		})
		return

	}
	if exist {
		log.Println("User Already exist.Please login")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User already exist.Please login",
		})

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error in Hashing the Password")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in Hashing the password",
		})
		return
	}

	user.Password = string(hashedPassword)
	user.Role = 0
	user.Timestamp = time.Now().UTC()

	_, err = DBHelper.InsertIntoUser(&user)
	if err != nil {
		log.Println("Error in Inserting to Users collection ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in inserting to users collection",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user inserted successfully",
		"user":    user,
	})

}
