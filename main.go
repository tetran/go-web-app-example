package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	// No need to set `Addr` field, because we use the listener passed as an argument
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Start the server in a goroutine
	eg.Go(func() error {
		// http.ErrServerClosed means the server was closed gracefully, so it's not an error.
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}

		return nil
	})

	// Wait for the context to be done
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Wait for the goroutine to finish
	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		log.Printf("Specify the port number as the first argument")
		os.Exit(1)
	}
	p := os.Args[1]

	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	if err := run(context.Background(), l); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}
