package auth

import (
	"log/slog"
	"net/http"
	jwt_encoder "url-shortener-api/internal/lib/jwt"
	"url-shortener-api/internal/lib/logger/sl"
)

func New(log *slog.Logger, encoder *jwt_encoder.JWTEncoder) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(slog.String("component", "middleware/auth"))

		log.Info("middleware auth enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				log.Error("missing token header")
				return
			}
			tokenString = tokenString[len("Bearer "):]

			err := encoder.VerifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				log.Error("invalid token", sl.Err(err))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
