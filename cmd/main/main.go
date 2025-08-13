package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"url-shorter/internal/http-server/handlers"
	del "url-shorter/internal/http-server/handlers/delete"
	"url-shorter/internal/http-server/handlers/url/save"

	_ "golang.org/x/crypto/bcrypt"
	"url-shorter/internal/config"
	"url-shorter/internal/http-server/handlers/redirect"
	authmw "url-shorter/internal/http-server/middleware"
	"url-shorter/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting url-shorter",
		slog.String("env", cfg.Env),
		slog.String("version", "123"))
	log.Debug("debug logging enabled")
	log.Error("error message enabled")

	store, err := storage.New(cfg.Storage)
	if err != nil {
		log.Error("failed to connect to storage", "error", err)
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// 1) Публичный логин: отдаём access/refresh
	router.Post("/auth/login", handlers.Login(store.DB))

	// 2) Приватные URL‑эндпоинты: доступ только с JWT
	router.Route("/url", func(r chi.Router) {
		r.Use(authmw.Auth) // наш JWT‑middleware

		r.Post("/", save.New(log, store))
		r.Delete("/{alias}", del.New(log, store))
		r.Get("/{alias}", redirect.New(log, store))

		r.Get("/{alias}/info", func(w http.ResponseWriter, r *http.Request) {
			alias := chi.URLParam(r, "alias")
			data, err := store.GetByAlias(alias)
			if err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
		})
	})

	// health
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "router not matched", http.StatusNotFound)
	})
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	log.Info("твой первый сервер запущен ты ПЗДЦ молодчага", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
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
