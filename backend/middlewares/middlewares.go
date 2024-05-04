package middlewares

import (
	"ecommerce/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var DBHelper *database.Mongo

// Authenticate will work as middleware to validate the JWT token
func Authenticate(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tokenString, err := c.Cookie("Authorization")
	secret := os.Getenv("SECRET_KEY")

	if err != nil {
		log.Println("Error in Getting the Tokenstring from cookie ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Unable to get the token string",
		})
		return

	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please login",
			})
			return

		}
		email := claims["sub"].(string)
		if len(email) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please login",
			})
			return
		}
		c.Set("user", email)
		c.Next()

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Please login",
		})
		return

	}

}

// IsAdmin middleware will check if the user is admin or not
func IsAdmin(c *gin.Context) {
	email, exist := c.Get("user")
	if !exist {
		log.Println("Please login")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Please login",
		})
		return

	}

	exist, err, user := DBHelper.CheckUser(email.(string))
	if !exist || err != nil {
		log.Println("Error in getting the user for admin check ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "SOme error occured in getting user",
		})
		return
	}
	if user.Role == 1 {
		log.Println("User is not Admin")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User is not Admin.",
		})
		return
	}
	c.Next()
}
