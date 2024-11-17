package http

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tender_management/internal/usecase/token"
)

type casbinPermission struct {
	enforcer *casbin.Enforcer
}

func (c *casbinPermission) GetRole(ctx *gin.Context) (string, int) {
	Token := ctx.GetHeader("Authorization")
	if Token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
		return "Unauthorized", http.StatusUnauthorized
	}
	claims, err := token.ExtractClaims(Token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return "Unauthorized", http.StatusUnauthorized
	}
	role, ok := claims["role"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role is empty"})
		return "Unauthorized", http.StatusUnauthorized
	}
	ctx.Set("user_id", claims["user_id"])
	return role, http.StatusOK
}

func (c *casbinPermission) CheckPermission(ctx *gin.Context) (bool, error) {
	subject, status := c.GetRole(ctx)
	if status != http.StatusOK {
		return false, errors.New("error while getting a role: " + subject)
	}
	acrtion := ctx.Request.Method
	object := ctx.FullPath()
	fmt.Println("subject", subject, "action", acrtion, "object", object)

	allow, err := c.enforcer.Enforce(subject, object, acrtion)
	if err != nil {
		return false, err
	}
	return allow, nil
}

func PermissionMiddleware(enf *casbin.Enforcer) gin.HandlerFunc {
	casbHandler := &casbinPermission{enforcer: enf}
	return func(ctx *gin.Context) {
		res, err := casbHandler.CheckPermission(ctx)

		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}
		if !res {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "You dont have permission"})
			return
		}
		auth := ctx.GetHeader("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "You dont have permission"})
			return
		}

		valid, err := token.ValidateToken(auth)
		if err != nil || !valid {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Token invalid: %s", err)})
			return
		}

		claims, err := token.ExtractClaims(auth)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Token invalid claims: %s", err)})
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Cors middleware triggered")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
