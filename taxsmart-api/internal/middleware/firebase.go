package middleware

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/taxsmart/taxsmart-api/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userID"

// FirebaseAuth creates a middleware that verifies Firebase ID tokens
func FirebaseAuth(app *firebase.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				response.Unauthorized(w, "Invalid authorization header format")
				return
			}

			client, err := app.Auth(r.Context())
			if err != nil {
				response.InternalError(w, "Failed to initialize auth client")
				return
			}

			token, err := client.VerifyIDToken(r.Context(), tokenString)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the user ID from the context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
