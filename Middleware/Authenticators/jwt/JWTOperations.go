package jwt

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define a secret key (store securely in environment variables)
var jwtSecret = []byte("your-secure-secret")

// Claims struct for JWT payload
type Claims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT for a client
func GenerateJWT(clientID string) (string, error) {
	claims := &Claims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24-hour expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a JWT and extracts claims
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("Authorization Header:", authHeader)

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := ValidateJWT(token)
		if err != nil {
			fmt.Println("Error validating token:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add client ID to context (optional)
		r = r.WithContext(context.WithValue(r.Context(), "client_id", claims.ClientID))
		fmt.Println("Client ID from token:", claims.ClientID)

		next.ServeHTTP(w, r)
	})
}
