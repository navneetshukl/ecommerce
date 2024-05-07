package routes

import (
	"ecommerce/database"
	"ecommerce/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	exist, err, _ := DBHelper.CheckUser(user.Email)

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
	user.Role = 2
	user.Timestamp = time.Now().UTC()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

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

// Login

func Login(c *gin.Context) {
	var user map[string]interface{}

	err := c.BindJSON(&user)
	if err != nil {
		log.Println("Error in reading the body response ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error in reading the body",
		})

		return
	}

	eLen := len(user["email"].(string))
	pLen := len(user["password"].(string))

	if eLen == 0 || pLen == 0 {
		log.Println("Enter Email and Password")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "please enter the email and password",
		})
		return
	}

	exists, err, data := DBHelper.CheckUser(user["email"].(string))

	if err != nil {
		log.Println("Error in checking the user from database ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error in getting the user information",
		})
		return
	}

	if !exists {
		log.Println("User does not exist")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User does not exist",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user["password"].(string)))
	if err != nil {

		log.Println("Password does not match")
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Password does not match",
		})

		return

	}

	// Create JWT token

	secret := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": data.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Println("Error in signing the token ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Some Error Occured.Please try again",
		})
		return
	}

	// Save this JWT token to Cookie

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, int(time.Hour*24*30), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User Login Successfully",
		"user": gin.H{
			"name":  data.Name,
			"email": data.Email,
			"phone": data.Phone,
			"token": tokenString,
		},
	})

}

func Test(c *gin.Context) {

	email, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in",
		"Value":   email.(string),
		"role":    "Admin",
	})
}
