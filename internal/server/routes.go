// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"issuer-service-go/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *FiberServer) RegisterRoutes(handler *Handler) {
	s.App.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	v1 := s.App.Group(config.GetConfig().ServerConfig.BasePath)
	v1.Get("/auth/*", notImplemented)
	v1.Get("/discovery/:realm", handler.DiscoveryHandler)
	v1.Get("/certs/:realm", handler.JwksHandler)
	v1.Get("/issuer/:realm", handler.IssuerHandler)

	auth := s.App.Group("/auth/realms/:realm")
	auth.Get("/protocol/openid-connect/auth/*", notImplemented)
	auth.Get("/.well-known/openid-configuration", handler.DiscoveryHandler)
	auth.Get("/protocol/openid-connect/certs", handler.JwksHandler)
	auth.Get("/", handler.IssuerHandler)
}

func notImplemented(c *fiber.Ctx) error {
	log.Debug().Msg("Request received on not implemented endpoint")
	return c.Status(fiber.StatusNotImplemented).SendString("This endpoint is unfortunately not available for Stargate")
}
