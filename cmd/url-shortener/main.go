package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener-api/internal/config"
	"url-shortener-api/internal/http-server/handlers/url/redirect"
	"url-shortener-api/internal/http-server/handlers/url/remove"
	urlSave "url-shortener-api/internal/http-server/handlers/url/save"
	"url-shortener-api/internal/http-server/handlers/user/login"
	userSave "url-shortener-api/internal/http-server/handlers/user/save"
	"url-shortener-api/internal/http-server/middleware/auth"
	mwLogger "url-shortener-api/internal/http-server/middleware/logger"
	jwt_encoder "url-shortener-api/internal/lib/jwt"
	"url-shortener-api/internal/lib/logger/sl"
	"url-shortener-api/internal/storage/mongodb"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("starting url-shortener api")
	logger.Debug("debug message")

	storage, err := mongodb.New(cfg.StoragePath)
	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	defer func(storage *mongodb.Storage) {
		err := storage.Disconnect()
		if err != nil {
			logger.Error("failed to disconnect", sl.Err(err))
			os.Exit(1)
		}
	}(storage)

	jwtEncoder, err := jwt_encoder.New(cfg.SecretKey)
	if err != nil {
		logger.Error("failed to create jwt encoder", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(auth.New(logger, jwtEncoder))

		r.Post("/", urlSave.New(logger, storage))
		r.Delete("/{alias}", remove.New(logger, storage))
	})

	router.Get("/{alias}", redirect.New(logger, storage))
	router.Post("/register", userSave.New(logger, storage))
	router.Post("/login", login.New(logger, storage, jwtEncoder))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
