package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SuleymanyanArkadi/eventhub/internal/logging"
	"github.com/SuleymanyanArkadi/eventhub/internal/reqid"
	"github.com/SuleymanyanArkadi/eventhub/internal/store"
)

func main() {
	// Создаём in-memory store
	s := store.NewMemoryStore()

	// Создаём маршруты, связанные с store
	mux := makeHandlers(s)

	// Оборачиваем middleware: сначала reqid (чтобы id был в контексте), затем logging
	handler := reqid.Middleware(logging.Middleware(mux))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
		close(idleConnsClosed)
	}()

	log.Printf("API server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	<-idleConnsClosed
	log.Println("server stopped")
}
