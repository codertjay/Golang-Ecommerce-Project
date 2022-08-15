package middleware

import (
	"Golang-Ecommerce-Project/tokens"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(400, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}
		claims, err := tokens.ValidateToken(clientToken)
		if err != "" {
			c.JSON(500, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
