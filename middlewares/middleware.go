package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("bajSBHJ344A35sfga0vaa!#44") // Ganti dengan kunci rahasia Anda

// Claims defines the structure for JWT claims
type Claims struct {
	UserID int64  `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func JWTMiddleware(ctx *context.Context) {

	// Ambil token dari header Authorization
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Authorization header missing"}, false, false)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token format"}, false, false)
		return
	}

	// Parsing token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.JSON(map[string]string{"error": "Invalid token"}, false, false)
		return
	}

	// Check user role and restrict access accordingly
	requestPath := ctx.Input.URL()
	if strings.HasPrefix(requestPath, "/admin/") && claims.Role != "admin" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "Access denied for non-admin users"}, false, false)
		return
	}

	if strings.HasPrefix(requestPath, "/user/") && claims.Role != "user" {
		ctx.Output.SetStatus(http.StatusForbidden)
		ctx.Output.JSON(map[string]string{"error": "Access denied for non-user users"}, false, false)
		return
	}

	// Token valid and role authorized, add user info to context
	ctx.Input.SetData("id", claims.UserID)
	ctx.Input.SetData("email", claims.Email)
	ctx.Input.SetData("role", claims.Role)

	// Log klaim yang berhasil dipetakan
	fmt.Printf("JWT Claims: UserID=%d, Email=%s, Role=%s\n", claims.UserID, claims.Email, claims.Role)
}
