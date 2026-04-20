package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/model"
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
		if !model.IsAdminRole(claims.Role) {
			support.RespondError(c, support.NewError(http.StatusForbidden, "forbidden", "admin access required"))
			c.Abort()
			return
		}
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := MustClaims(c)
		if claims.UserID == 0 {
			support.RespondError(c, support.NewError(http.StatusUnauthorized, "unauthorized", "login required"))
			c.Abort()
			return
		}
		if claims.Role != model.RoleSuperAdmin && !slices.Contains(roles, claims.Role) {
			support.RespondError(c, support.NewError(http.StatusForbidden, "forbidden", "insufficient role permissions"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func PermissionRequired(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := MustClaims(c)
		if claims.UserID == 0 {
			support.RespondError(c, support.NewError(http.StatusUnauthorized, "unauthorized", "login required"))
			c.Abort()
			return
		}
		if !model.HasAllPermissions(claims.Role, permissions...) {
			support.RespondError(c, support.NewError(http.StatusForbidden, "forbidden", "insufficient permission"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func MustClaims(c *gin.Context) support.TokenClaims {
	raw, _ := c.Get(contextClaimsKey)
	if claims, ok := raw.(support.TokenClaims); ok {
		return claims
	}
	return support.TokenClaims{}
}
