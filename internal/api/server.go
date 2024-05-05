package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-fitness/external/config"
	"go-fitness/external/logger/sl"
	"go.uber.org/fx"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const EnvPath = ".env"

func NewServer(cfg *config.Config, r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    cfg.HTTPServer.ApiPort,
		Handler: r,
		/*ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,*/
		WriteTimeout: time.Minute * 10,
		ReadTimeout:  time.Minute * 10,
		IdleTimeout:  time.Minute * 10,
	}
}

func RunServer(
	lc fx.Lifecycle,
	log *slog.Logger,
	server *http.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Starting server", slog.Any("config", server.Addr))

				if os.Getenv("JWT_SECRET") == "" {
					appToken := generateAppToken()
					log.Info("Generated app token", slog.String("token", appToken))
					saveTokenToENV(appToken)
				}

				if err := server.ListenAndServe(); err != nil {
					log.Error("Server failed", sl.Err(err))
				} else {
					log.Info("Server started")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Error("Server stopped")
			return server.Shutdown(ctx)
		},
	})
}

func generateAppToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	token := hex.EncodeToString(b)
	return token
}

func saveTokenToENV(token string) {
	file, err := os.OpenFile(EnvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open env file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("\nJWT_SECRET=%s", token)); err != nil {
		log.Fatalf("Failed to write to env file: %v", err)
	}
}
