package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/devanadindra/signlink-mobile/back-end/utils/config"
	"github.com/devanadindra/signlink-mobile/back-end/utils/logger"
	wireinject "github.com/devanadindra/signlink-mobile/back-end/wire"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
)

func main() {
	conf := config.NewConfig()

	err := validator.New().Struct(conf)
	if err != nil {
		panic(err)
	}

	dependency, err := wireinject.InitializeDependency(conf)
	if err != nil {
		panic(err)
	}
	defer dependency.Close()

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://10.0.2.2:8080",
			"http://10.0.3.2:8080",
			"http://192.168.1.10:8080",
		},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Frontend"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	})

	handler := dependency.GetHandler()

	// Apply CORS middleware
	handlerWithCORS := c.Handler(handler)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: handlerWithCORS,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Error listen and serve : %v\n", err))
		}
	}()

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()

	logger.Trace(ctx, "shutting down gracefully, please wait a moment ...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Trace(ctx, "Server gracefully shutdown")
}
