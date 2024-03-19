package login

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shortener-api/internal/lib/api/response"
	"url-shortener-api/internal/lib/hasher"
	jwt_encoder "url-shortener-api/internal/lib/jwt"
	"url-shortener-api/internal/lib/logger/sl"
	"url-shortener-api/internal/lib/user"
)

type Request struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Response    resp.Response
	AccessToken string `json:"access_token"`
}

type UserGetter interface {
	GetUser(email string) (user.User, error)
}

func New(log *slog.Logger, userGetter UserGetter, jwtEncoder *jwt_encoder.JWTEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode json", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed", sl.Err(err))

			render.JSON(w, r, resp.Error("validation failed"))
			return
		}

		user, err := userGetter.GetUser(req.Email)
		if err != nil {
			log.Error("failed to login", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to login"))
			return
		}

		reqPassword := hasher.GetMD5Hash(req.Password)
		if reqPassword != user.Password {
			log.Error("wrong password")

			render.JSON(w, r, resp.Error("wrong password"))
			return
		}

		tokenString, err := jwtEncoder.CreateToken(&user)
		if err != nil {
			log.Error("user login error", sl.Err(err))

			render.JSON(w, r, resp.Error("user login error"))
			return
		}

		log.Info("login successful", slog.String("token", tokenString))

		render.JSON(w, r, Response{
			Response:    resp.Ok(),
			AccessToken: tokenString,
		})
	}
}
