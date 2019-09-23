package server

import (
	// Native packages
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Third parties
	"github.com/go-chi/chi"
)

func serveWithGracefulShutdown(server *http.Server, timeout time.Duration, proc func() error) error {
	// Create listener for the 'SIGTERM'
	// from kernel
	trigger := make(chan os.Signal, 1)
	signal.Notify(trigger, os.Interrupt, syscall.SIGTERM)

	// Wait for 'SIGTERM' from kernel
	var errRunning error
	go func() { errRunning = proc() }()
	<-trigger

	// Create the cancelable context for
	// help cancel the halted shutdown
	// process
	srvCtx, srvCancel := context.WithTimeout(context.Background(), timeout)
	defer srvCancel()

	// Perform shutdown then wait until
	// the server finished the shutdown
	// process or the timeout had been
	// reached
	errShutdown := server.Shutdown(srvCtx)
	if errShutdown == nil {
		return errRunning
	}

	return errShutdown
}

// ServeHTTP serves the server
// on HTTP protocol.
func ServeHTTP(routers *chi.Mux, address string, timeout time.Duration) error {
	server := &http.Server{
		Handler: routers,
		Addr:    address,
	}

	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServe()
	})
}

// ServeHTTPS serves the server
// on HTTPS protocol.
func ServeHTTPS(routers *chi.Mux, publicKey, privateKey, address string, timeout time.Duration) error {

	server := &http.Server{
		Handler: routers,
		Addr:    address,
	}

	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServeTLS(publicKey, privateKey)
	})
}

// End-of-file
