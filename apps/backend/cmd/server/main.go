package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	"github.com/go-chi/httprate"

	"github.com/react-go-monorepo/backend/internal/config"
	"github.com/react-go-monorepo/backend/internal/db"
	"github.com/react-go-monorepo/backend/internal/handler"
)

func main() {
	// 設定の読み込み
	cfg := config.Load()

	// ロガーの初期化
	logger := httplog.NewLogger("react-go-monorepo", httplog.Options{
		JSON:    cfg.IsProduction(),
		Concise: cfg.IsProduction(),
	})

	// データベース接続
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	// ルーターの設定
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// レート制限
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// ハンドラーの設定
	h := handler.New(database, cfg)

	// ルート
	r.Route("/api/v1", func(r chi.Router) {
		// ヘルスチェック
		r.Get("/health", h.HealthCheck)

		// 認証
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", h.Login)
			r.Post("/refresh", h.RefreshToken)
		})

		// 認証が必要なルート
		r.Group(func(r chi.Router) {
			r.Use(h.AuthMiddleware)

			// Todos
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", h.GetTodos)
				r.Post("/", h.CreateTodo)
				r.Route("/{todoId}", func(r chi.Router) {
					r.Get("/", h.GetTodoById)
					r.Put("/", h.UpdateTodo)
					r.Delete("/", h.DeleteTodo)
				})
			})
		})
	})

	// サーバーの起動
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info().Str("port", cfg.Port).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("Server failed to start")
			os.Exit(1)
		}
	}()

	<-done
	logger.Info().Msg("Server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
		os.Exit(1)
	}

	logger.Info().Msg("Server stopped")
}
