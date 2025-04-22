package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/harshgupta9473/assignment_makerable/utils"
)

type UserClaims struct {
	UserID     int64    `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`        // "doctor" or "patient"
	IsVerified bool   `json:"is_verified"` // email verification
	IsApproved bool   `json:"is_approved"` // only for doctors
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64, email, role string, isVerified, isApproved bool) (string, error) {
	claims := UserClaims{
		UserID:     userID,
		Email:      email,
		Role:       role,
		IsVerified: isVerified,
		IsApproved: isApproved,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.JwtSecretKey))
}

func ValidateJWT(tokenStr string, jwtSecretKey string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken == "" || !strings.HasPrefix(authToken, "Bearer ") {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "No auth token provided",
				Error:   "INVALID_AUTH_TOKEN",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authToken, "Bearer ")
		claims, err := ValidateJWT(tokenStr,utils.JwtSecretKey)
		if err != nil {
			utils.WriteJson(w, http.StatusUnauthorized, utils.APIResponse{
				Status:  "error",
				Message: "LogIN Again",
				Error:   "INVALID_JWT",
			})
			return
		}
		
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserContext(r *http.Request)(*UserClaims,error){
	claim,ok:=r.Context().Value("user").(*UserClaims)
	if !ok{
		return nil,fmt.Errorf("error extracting the user context")
	}
	return claim,nil
}