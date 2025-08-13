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

type JwksResponse struct {
	Keys []*jwks.Jwk `json:"keys"`
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

func NewDiscoveryInfo(issuerURL string, realm string) Discovery {
	return Discovery{
		IssuerURL: fmt.Sprintf("%s/auth/realms/%s", issuerURL, realm),
		JwksURL: fmt.Sprintf(
			"%s/auth/realms/%s/protocol/openid-connect/certs",
			issuerURL,
			realm,
		),
		AuthorizationEndpointURL: fmt.Sprintf(
			"%s/auth/realms/%s/protocol/openid-connect/auth",
			issuerURL,
			realm,
		),
		ResponseTypesSupported:           []string{"none"},
		SubjectTypesSupported:            []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
	}
}

func (h *Handler) DiscoveryHandler(c *fiber.Ctx) error {
	log.Debug().Msg("Request received on discovery endpoint")
	realm := c.Params("realm")

	log.Debug().Msgf("Request with following headers: %+v", c.GetReqHeaders())
	host := c.Get("X-Forwarded-Host")
	if host == "" {
		log.Error().Msg("X-Forwarded-Host header must be set in the request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "X-Forwarded-Host header must be set in the request",
		})
	}
	issuerURL := "https://" + host + config.GetConfig().PathPrefix

	return c.Status(fiber.StatusOK).JSON(NewDiscoveryInfo(issuerURL, realm))
}

func (h *Handler) JwksHandler(c *fiber.Ctx) error {
	realm := c.Params("realm")
	log.Debug().Msgf("Request received on certs endpoint for realm %s", realm)

	info := h.jwksProvider.GetJwks()
	response := &JwksResponse{
		Keys: info,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) IssuerHandler(c *fiber.Ctx) error {
	realm := c.Params("realm")
	log.Debug().Msgf("Request received on issuer endpoint for realm %s", realm)

	defaultRealm := h.jwksProvider.GetDefaultRealm(realm)

	return c.Status(fiber.StatusOK).JSON(defaultRealm)
}
