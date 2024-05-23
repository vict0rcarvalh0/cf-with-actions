package middlewares

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
)

type contextKey string

const UserContextKey contextKey = "userID"

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		// Removes the "Bearer " prefix
		if len(tokenString) < 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized", "message": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Supondo que o ID do usuário está armazenado na claim "userID"
			userID, err := claims["id"].(string)
			if !err {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Unauthorized", "message": "Invalid user ID"})
				return
			}

			// Adiciona o ID do usuário ao contexto da requisição
			ctx := context.WithValue(r.Context(), UserContextKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		} else {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Unauthorized", "message": "Invalid token"})
			return
		}
	})
}
