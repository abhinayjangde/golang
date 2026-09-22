package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/abhinayjangde/ecom/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: also inject logger dependency to log the errors in the middleware
func AuthMiddleware(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the presence of the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.CodeUnauthenticated)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.CodeUnauthenticated)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.CodeUnauthenticated)
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			httpx.Error(w, http.StatusUnauthorized, "Unauthorized", httpx.CodeUnauthenticated)
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string)
		email := claims["email"].(string)
		role := claims["role"].(string)

		// store data in context for further use in the request lifecycle
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) string {
	return ctx.Value("user_id").(string)
}

func EmailFromContext(ctx context.Context) string {
	return ctx.Value("email").(string)
}

func RoleFromContext(ctx context.Context) string {
	return ctx.Value("role").(string)
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if RoleFromContext(ctx) != "admin" {
			httpx.Error(
				w,
				http.StatusForbidden,
				"admin access required",
				httpx.CodeForbidden,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
