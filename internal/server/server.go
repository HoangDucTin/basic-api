package server

import (
	"github.com/HoangDucTin/basic-api/internal/router"
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

func ServeHttp(address string, timeout time.Duration) error {

	server := &http.Server{
		Handler: router.NewRouter(),
		Addr:    address,
	}

	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServe()
	})
}

func ServeHttps(publicKey, privateKey, address string, timeout time.Duration) error {

	server := &http.Server{
		Handler: router.NewRouter(),
		Addr:    address,
	}

	return serveWithGracefulShutdown(server, timeout, func() error {
		return server.ListenAndServeTLS(publicKey, privateKey)
	})
}

// End-of-file
