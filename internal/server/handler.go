// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"issuer-service-go/internal/config"
	"issuer-service-go/internal/jwks"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type HandlerInterface interface {
	DiscoveryHandler(c *fiber.Ctx) error
	JwksHandler(c *fiber.Ctx) error
	IssuerHandler(c *fiber.Ctx) error
}

type Handler struct {
	jwksProvider jwks.Provider
}

func NewHandler(jwksProvider jwks.Provider) *Handler {
	return &Handler{
		jwksProvider: jwksProvider,
	}
}

type Discovery struct {
	IssuerURL                        string   `json:"issuer"`
	JwksURL                          string   `json:"jwks_uri"`
	AuthorizationEndpointURL         string   `json:"authorization_endpoint"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

//nolint:golines // ignore this linter error
func NewDiscoveryInfo(issuerURL string, realm string) Discovery {
	return Discovery{
		IssuerURL:                        fmt.Sprintf("%s/auth/realms/%s", issuerURL, realm),
		JwksURL:                          fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/certs", issuerURL, realm),
		AuthorizationEndpointURL:         fmt.Sprintf("%s/auth/realms/%s/protocol/openid-connect/auth", issuerURL, realm),
		ResponseTypesSupported:           []string{"none"},
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}
}

func (h *Handler) DiscoveryHandler(c *fiber.Ctx) error {
	log.Debug().Msg("Request received on discovery endpoint")
	realm := c.Params("realm")

	return c.Status(fiber.StatusOK).JSON(NewDiscoveryInfo(config.GetConfig().IssuerURL, realm))
}

func (h *Handler) JwksHandler(c *fiber.Ctx) error {
	realm := c.Params("realm")
	log.Debug().Msgf("Request received on certs endpoint for realm %s", realm)

	info := h.jwksProvider.GetJwks()

	return c.Status(fiber.StatusOK).JSON(info)
}

func (h *Handler) IssuerHandler(c *fiber.Ctx) error {
	realm := c.Params("realm")
	log.Debug().Msgf("Request received on issuer endpoint for realm %s", realm)

	defaultRealm := h.jwksProvider.GetDefaultRealm(realm)

	return c.Status(fiber.StatusOK).JSON(defaultRealm)
}
