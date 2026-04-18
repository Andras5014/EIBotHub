package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/support"
)

const contextClaimsKey = "claims"

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func AuthRequired(tokens *support.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			support.RespondError(c, support.NewError(http.StatusUnauthorized, "unauthorized", "missing bearer token"))
			c.Abort()
			return
		}

		claims, err := tokens.Parse(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			support.RespondError(c, err)
			c.Abort()
			return
		}

		c.Set(contextClaimsKey, claims)
		c.Next()
	}
}

func AdminRequired(tokens *support.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthRequired(tokens)(c)
		if c.IsAborted() {
			return
		}

		claims := MustClaims(c)
		if claims.Role != "admin" {
			support.RespondError(c, support.NewError(http.StatusForbidden, "forbidden", "admin access required"))
			c.Abort()
			return
		}
	}
}

func MustClaims(c *gin.Context) support.TokenClaims {
	raw, _ := c.Get(contextClaimsKey)
	if claims, ok := raw.(support.TokenClaims); ok {
		return claims
	}
	return support.TokenClaims{}
}
