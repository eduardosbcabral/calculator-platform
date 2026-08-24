package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/eduardosbcabral/calculator-platform/internal/calculator"
)

func run(ctx context.Context) error {
	handler, err := newHandler(os.Getenv("STATIC_DIR"))
	if err != nil {
		return err
	}

	server := newServer(handler, cmp.Or(os.Getenv("PORT"), "8080"))
	return serve(ctx, server)
}

func newHandler(staticDir string) (http.Handler, error) {
	mux := http.NewServeMux()
	calculator.Register(mux)
	mux.HandleFunc("/healthz", health)

	if staticDir == "" {
		mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "Frontend is not configured.", http.StatusNotFound)
		})
		return mux, nil
	}

	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err != nil {
		return nil, fmt.Errorf("static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	return mux, nil
}

func newServer(handler http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func serve(ctx context.Context, server *http.Server) error {
	serverError := make(chan error, 1)
	go func() {
		slog.Info("server started", "address", server.Addr)
		serverError <- server.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
