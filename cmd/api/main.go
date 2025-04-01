// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"issuer-service-go/internal/config"
	"issuer-service-go/internal/jwks"
	"issuer-service-go/internal/server"
	"issuer-service-go/internal/version"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Info().Msg("shutting down gracefully...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), config.GetConfig().GracefulShutdownTimeout)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Info().Msg("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	log.Info().Msgf("%s\n", version.GetVersionInfo())

	appConfig := config.GetConfig()

	jwksProvider, err := jwks.NewFileProvider(&appConfig.JwksConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create JWKS file provider")
	}
	handler := server.NewHandler(jwksProvider)

	srv := server.New()
	srv.RegisterRoutes(handler)

	done := make(chan bool, 1)

	go func() {
		srvError := srv.Listen(fmt.Sprintf(":%d", appConfig.ServerConfig.Port))
		if srvError != nil {
			panic(fmt.Sprintf("http server error: %s", srvError))
		}
	}()

	go gracefulShutdown(srv, done)

	<-done
	log.Info().Msg("Graceful shutdown complete.")
}
