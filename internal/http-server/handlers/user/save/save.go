package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url-shortener-api/internal/lib/api/response"
	"url-shortener-api/internal/lib/logger/sl"
)

type Request struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type Response struct {
	Response resp.Response
	Id       string `json:"id"`
}

type UserSaver interface {
	CreateUser(email, password string) (string, error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
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

		if req.Password != req.ConfirmPassword {
			log.Error("password not equals confirm password")

			render.JSON(w, r, resp.Error("password not equals confirm password"))
			return
		}

		id, err := userSaver.CreateUser(req.Email, req.Password)
		if err != nil {
			log.Error("save error", sl.Err(err))

			render.JSON(w, r, resp.Error("save error"))
			return
		}

		log.Info("created successful", slog.String("id", id))

		render.JSON(w, r, Response{
			Response: resp.Ok(),
			Id:       id,
		})
	}
}
