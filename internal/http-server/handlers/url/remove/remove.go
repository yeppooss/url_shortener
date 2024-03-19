package remove

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener-api/internal/lib/api/response"
	"url-shortener-api/internal/lib/logger/sl"
)

type URLRemover interface {
	DeleteUrl(alias string) error
}

func New(log *slog.Logger, urlSaver URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.Delete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias param is empty")

			render.JSON(w, r, resp.Error("alias param is empty"))
			return
		}
		err := urlSaver.DeleteUrl(alias)
		if err != nil {
			log.Error("failed to delete url ", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete url"))
			return
		}

		render.JSON(w, r, nil)
	}
}
