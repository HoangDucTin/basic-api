package server

import (
	// Native packages
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
func ServeHTTP(server *http.Server, timeout time.Duration) error {
	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServe()
	})
}

// ServeHTTPS serves the server
// on HTTPS protocol.
func ServeHTTPS(server *http.Server, publicKey, privateKey string, timeout time.Duration) error {
	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServeTLS(publicKey, privateKey)
	})
}

// Serve serves the server
// on both HTTP and HTTPS
// protocols.
func Serve(server *http.Server, publicKey, privateKey string, timeout time.Duration) error {
	trigger := make(chan error, 1)
	go func() { trigger <- ServeHTTP(server, timeout) }()
	go func() { trigger <- ServeHTTPS(server, privateKey, publicKey, timeout) }()

	if errServing := <-trigger; errServing != nil {
		srvCtx, srvCancel := context.WithTimeout(context.Background(), timeout)
		defer srvCancel()

		if err := server.Shutdown(srvCtx); err != nil {
			return err
		}
	}

	return nil
}

// End-of-file
