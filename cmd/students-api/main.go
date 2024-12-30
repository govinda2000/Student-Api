package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/govinda2000/students-api/internal/config"
	"github.com/govinda2000/students-api/internal/http/handlers/student"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started on port:", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	// TO-DO: learn about these signals
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("faied to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
