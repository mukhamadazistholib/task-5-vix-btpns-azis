package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"task-vix-btpns/app/auth"
)

//function to protect routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") //Get bearer token
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}

		err := auth.ValidateToken(strings.Split(tokenString, "Bearer ")[1]) //Validate token
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
