// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"issuer-service-go/internal/config"
	"issuer-service-go/internal/jwks"
	"issuer-service-go/internal/server"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup code here
	err := os.Setenv("CERT_MOUNT_PATH", "./internal/api/router_testdata")
	if err != nil {
		fmt.Printf("failed to set environment variable 'CERT_MOUNT_PATH': %v\n", err)
	}
	err = os.Setenv("CERT_FILE_NEXT", "next-tls.crt")
	if err != nil {
		fmt.Printf("failed to set environment variable 'CERT_FILE_NEXT': %v\n", err)
	}
	err = os.Setenv("KID_FILE_NEXT", "next-tls.kid")
	if err != nil {
		fmt.Printf("failed to set environment variable 'KID_FILE_NEXT': %v\n", err)
	}
	err = os.Setenv("CERT_FILE_ACTIVE", "tls.crt")
	if err != nil {
		fmt.Printf("failed to set environment variable 'CERT_FILE_ACTIVE': %v\n", err)
	}
	err = os.Setenv("KID_FILE_ACTIVE", "tls.kid")
	if err != nil {
		fmt.Printf("failed to set environment variable 'KID_FILE_ACTIVE': %v\n", err)
	}
	err = os.Setenv("CERT_FILE_PREV", "prev-tls.crt")
	if err != nil {
		fmt.Printf("failed to set environment variable 'CERT_FILE_PREV': %v\n", err)
	}
	err = os.Setenv("KID_FILE_PREV", "prev-tls.kid")
	if err != nil {
		fmt.Printf("failed to set environment variable 'KID_FILE_PREV': %v\n", err)
	}

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestHealthRoute(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Test /health endpoint",
			route:        "/health",
			expectedCode: 200,
		},
	}

	srv := server.New()
	handler := server.NewHandler(nil) // jwksProvider is not needed for that test
	srv.RegisterRoutes(handler)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			resp, _ := srv.Test(req, 1)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		})
	}
}

func TestAuthRoute(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Test /auth endpoint #1",
			route:        config.GetConfig().ServerConfig.BasePath + "/auth",
			expectedCode: 501,
		},
		{
			description:  "Test /auth endpoint #2",
			route:        config.GetConfig().ServerConfig.BasePath + "/auth/test",
			expectedCode: 501,
		},
		{
			description:  "Test /auth endpoint #3",
			route:        config.GetConfig().ServerConfig.BasePath + "/auth/test/test",
			expectedCode: 501,
		},
		{
			description:  "Test /auth endpoint #4",
			route:        "/auth/realms/default/protocol/openid-connect/auth",
			expectedCode: 501,
		},
		{
			description:  "Test /auth endpoint #5",
			route:        "/auth/realms/default/protocol/openid-connect/auth/test",
			expectedCode: 501,
		},
		{
			description:  "Test /auth endpoint #6",
			route:        "/auth/realms/default/protocol/openid-connect/auth/test/test",
			expectedCode: 501,
		},
	}

	srv := server.New()
	handler := server.NewHandler(nil) // jwksProvider is not needed for that test
	srv.RegisterRoutes(handler)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			resp, _ := srv.Test(req, 1)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)
		})
	}
}

func TestDiscoveryRoute(t *testing.T) {
	issuerUrl := "localhost:8080"

	tests := []struct {
		description      string
		route            string
		pathPrefix       string
		headers          map[string]string
		expectedCode     int
		expectedResponse server.Discovery
		expectedError    fiber.Error
	}{
		{
			description:  "Test /discovery/default endpoint #1",
			route:        config.GetConfig().ServerConfig.BasePath + "/discovery/default",
			headers:      map[string]string{"X-Forwarded-Host": issuerUrl},
			expectedCode: 200,

			expectedResponse: server.Discovery{
				IssuerURL:                        "https://" + issuerUrl + "/auth/realms/default",
				JwksURL:                          "https://" + issuerUrl + "/auth/realms/default/protocol/openid-connect/certs",
				AuthorizationEndpointURL:         "https://" + issuerUrl + "/auth/realms/default/protocol/openid-connect/auth",
				ResponseTypesSupported:           []string{"none"},
				SubjectTypesSupported:            []string{"public"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"}},
		},
		{
			description:  "Test /discovery/other-realm endpoint #2",
			route:        config.GetConfig().ServerConfig.BasePath + "/discovery/other-realm",
			headers:      map[string]string{"X-Forwarded-Host": issuerUrl},
			expectedCode: 200,
			expectedResponse: server.Discovery{
				IssuerURL:                        "https://" + issuerUrl + "/auth/realms/other-realm",
				JwksURL:                          "https://" + issuerUrl + "/auth/realms/other-realm/protocol/openid-connect/certs",
				AuthorizationEndpointURL:         "https://" + issuerUrl + "/auth/realms/other-realm/protocol/openid-connect/auth",
				ResponseTypesSupported:           []string{"none"},
				SubjectTypesSupported:            []string{"public"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"}},
		},
		{
			description:  "Test /discovery/default endpoint #3",
			route:        "/auth/realms/default/.well-known/openid-configuration",
			headers:      map[string]string{"X-Forwarded-Host": issuerUrl},
			expectedCode: 200,
			expectedResponse: server.Discovery{
				IssuerURL:                        "https://" + issuerUrl + "/auth/realms/default",
				JwksURL:                          "https://" + issuerUrl + "/auth/realms/default/protocol/openid-connect/certs",
				AuthorizationEndpointURL:         "https://" + issuerUrl + "/auth/realms/default/protocol/openid-connect/auth",
				ResponseTypesSupported:           []string{"none"},
				SubjectTypesSupported:            []string{"public"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"}},
		},
		{
			description:  "Test /discovery/other-realm endpoint #4",
			route:        "/auth/realms/other-realm/.well-known/openid-configuration",
			headers:      map[string]string{"X-Forwarded-Host": issuerUrl},
			expectedCode: 200,
			expectedResponse: server.Discovery{
				IssuerURL:                        "https://" + issuerUrl + "/auth/realms/other-realm",
				JwksURL:                          "https://" + issuerUrl + "/auth/realms/other-realm/protocol/openid-connect/certs",
				AuthorizationEndpointURL:         "https://" + issuerUrl + "/auth/realms/other-realm/protocol/openid-connect/auth",
				ResponseTypesSupported:           []string{"none"},
				SubjectTypesSupported:            []string{"public"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"}},
		},
		{
			description:  "Test /discovery/default endpoint with no X-Forwarded-Host header #5",
			route:        config.GetConfig().ServerConfig.BasePath + "/discovery/default",
			headers:      map[string]string{},
			expectedCode: 400,
			expectedError: fiber.Error{
				Message: "X-Forwarded-Host header must be set in the request",
				Code:    400,
			},
		},
		{
			description:  "Test /discovery/default endpoint with no X-Forwarded-Host header value #6",
			route:        config.GetConfig().ServerConfig.BasePath + "/discovery/default",
			headers:      map[string]string{"X-Forwarded-Host": ""},
			expectedCode: 400,
			expectedError: fiber.Error{
				Message: "X-Forwarded-Host header must be set in the request",
				Code:    400,
			},
		},
		{
			description:  "Test /discovery/default endpoint for spacegate #7",
			route:        "/auth/realms/default/.well-known/openid-configuration",
			pathPrefix:   "/spacegate",
			headers:      map[string]string{"X-Forwarded-Host": issuerUrl},
			expectedCode: 200,
			expectedResponse: server.Discovery{
				IssuerURL:                        "https://" + issuerUrl + "/spacegate/auth/realms/default",
				JwksURL:                          "https://" + issuerUrl + "/spacegate/auth/realms/default/protocol/openid-connect/certs",
				AuthorizationEndpointURL:         "https://" + issuerUrl + "/spacegate/auth/realms/default/protocol/openid-connect/auth",
				ResponseTypesSupported:           []string{"none"},
				SubjectTypesSupported:            []string{"public"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"}},
		},
	}

	srv := server.New()
	handler := server.NewHandler(nil) // jwksProvider is not needed for that test
	srv.RegisterRoutes(handler)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			config.GetConfig().PathPrefix = tt.pathPrefix

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Set("X-Forwarded-Host", tt.headers["X-Forwarded-Host"])
			resp, _ := srv.Test(req, 1)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)

			var actualResponse server.Discovery
			err := json.NewDecoder(resp.Body).Decode(&actualResponse)
			if err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			assert.Equalf(t, tt.expectedResponse, actualResponse, tt.description)
		})
	}
}

//nolint:golines // ignore this linter error
func TestJwksRoute(t *testing.T) {
	jwkPrev := jwks.Jwk{
		Kid:     "5A9C11C2-A370-473D-AB2B-4B8BC247724C",
		Kty:     "RSA",
		Alg:     "RS256",
		Use:     "sig",
		N:       "4nuBr5l7UtS5X-aRMG5_XNtQDvz-NddyCdnAcbquMpx8WEHDRtj47kDmeb01tvrWqZNgrRyQnDl1xQ5qkV0BXoS7n_iUWtrgxJYfprIqHoEFFclugLQiyzbKkez4Y6gw0Zaz7bbB5FRbBKc4Md0DOmJXn-m-6smu5--6FkUixjXZi24YZSVIYjpDxiJpVJmVotaTrOX615VWolk9wdJ0d6dKfIim9YdMFPgJbiLsHL3wi64m8D8TqzXzJynwED4mAW-CKnPp9ueSsQLkVZLMAYmqGt8upsTe046j9y73BVxR_-YwJ7utOkiD2C3jGL_6ex4WBIhEedAC2dO4sBrELw",
		E:       "AQAB",
		X5c:     []string{"MIIC4TCCAcmgAwIBAgIUUXNbl9Vgby/oKY1Bqyz5nAFiUEkwDQYJKoZIhvcNAQELBQAwADAeFw0yNTA0MTExMzIyMTZaFw0yODAxMDYxMzIyMTZaMAAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDie4GvmXtS1Llf5pEwbn9c21AO/P4113IJ2cBxuq4ynHxYQcNG2PjuQOZ5vTW2+tapk2CtHJCcOXXFDmqRXQFehLuf+JRa2uDElh+msioegQUVyW6AtCLLNsqR7PhjqDDRlrPttsHkVFsEpzgx3QM6Ylef6b7qya7n77oWRSLGNdmLbhhlJUhiOkPGImlUmZWi1pOs5frXlVaiWT3B0nR3p0p8iKb1h0wU+AluIuwcvfCLribwPxOrNfMnKfAQPiYBb4Iqc+n255KxAuRVkswBiaoa3y6mxN7TjqP3LvcFXFH/5jAnu606SIPYLeMYv/p7HhYEiER50ALZ07iwGsQvAgMBAAGjUzBRMB0GA1UdDgQWBBRY7HmChVPSIqHSPwdEZ9qP5hW1pTAfBgNVHSMEGDAWgBRY7HmChVPSIqHSPwdEZ9qP5hW1pTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQA7wNOb93eMXpbu0GqTugeK9C1+4R3lfZMauuPNMdZk0ylhzvS9uRMcfre18hWJuWBykap+8vKIVs/Ia611rPP5ye+jzhJ4MUt/G8Jf9eEbgmF3wKKuCItI4tN0plLRPntqgPz/uTi/pwDovOV/meytmIP+hZ7kN3r2soOHhqtSVYbFpKmMpfu8mRGTCBMJt7Hv6d/tNFcAJkFb/Y0m48Eci5n5Doe9pjrRWOHVeJnnF3wSlETm7WNbgpkreXrG8wFx9/iZgirZ40WDiwHVPxen2piI+esCienCWb1clBmT1gPh6DC+zNHm7F1bbJWLKonFpCyO5gZegL+85uV4bDFG"},
		X5t:     "ClJICZjMGGZ2XbEHsPsqhNJq4b4",
		X5tS256: "MbjFiLW3TEVuUkIFDlQe8QSs0TQ_ofwSppNTtn-pIi0",
	}

	jwkActive := jwks.Jwk{
		Kid:     "F7959F8A-EC16-44BC-9F77-2A6F9580BDB4",
		Kty:     "RSA",
		Alg:     "RS256",
		Use:     "sig",
		N:       "y0skyMX46fzroq2Ma2pr1iP-Rt-x3IKufm6rf54vwcq_jxYPBajNREM0dtKfjj1p590Gme1-QQW3uS03eK5Rp5CNGGonFrzWsqlYa3dYgHpcZ6UgFGBPJvJCqBnEFP3d7zdg4GKOGDGv-KEM49bKm1qfIvxJ-JpATzv06vNptsGrtygol1rVbWkq8cFZ5mIzSe3Jk0vx8tw3rEint4uG8OHNWqfdHBKblTVjuW2w6cYr7gk6ujm9FswjkZ5us0mgBekw0prLK5bYwNzHERdFtvaCvOIwNZvqwsETQFpQFBwB_7kdEFfuSHbDeG0Mg5_aIikKom2TV-bEy21V6Sw_1Q",
		E:       "AQAB",
		X5c:     []string{"MIIC4TCCAcmgAwIBAgIUPmDMG0Hiqo2+DSBWFqRvk6p7SZswDQYJKoZIhvcNAQELBQAwADAeFw0yNTA0MTExMzI0MTRaFw0yODAxMDYxMzI0MTRaMAAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDLSyTIxfjp/OuirYxramvWI/5G37Hcgq5+bqt/ni/Byr+PFg8FqM1EQzR20p+OPWnn3QaZ7X5BBbe5LTd4rlGnkI0YaicWvNayqVhrd1iAelxnpSAUYE8m8kKoGcQU/d3vN2DgYo4YMa/4oQzj1sqbWp8i/En4mkBPO/Tq82m2wau3KCiXWtVtaSrxwVnmYjNJ7cmTS/Hy3DesSKe3i4bw4c1ap90cEpuVNWO5bbDpxivuCTq6Ob0WzCORnm6zSaAF6TDSmssrltjA3McRF0W29oK84jA1m+rCwRNAWlAUHAH/uR0QV+5IdsN4bQyDn9oiKQqibZNX5sTLbVXpLD/VAgMBAAGjUzBRMB0GA1UdDgQWBBThjgw4oMINmiRhVhjf1fKmFSkYkTAfBgNVHSMEGDAWgBThjgw4oMINmiRhVhjf1fKmFSkYkTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBMcu0NfhHMGV2JvUggDXoX5ws3RVtRS353n/JFtuF+ngH9zQgPb1NL/2MYe811Mp7mR399xdyfCHSglZ0uLsMIKZ6hyBCxvy9bLGrlLTQzC/azXyiegk+Q4vABXftSJ7x0a0H+zlDBdkYoda99Mcx4igJE3tpdeyTYpBZXhCcI6U6ZB0Ck/Wf/z36ncw7nEP0rCABV+cTtlJGPLhnU8AQpHvOSvSnRlC8hzSZX+n6+im7gRlegpBwtLy4KmThWYdcdnu/3o68iKeSJgykL/w16u3Uc/aYjBf2VwhJmNKmeRn2FMzY+kxr7MoVXV2a91mtM+t5WfvIK9jZL/11i2ZMG"},
		X5t:     "Zilf2ZAyRv8IUWg7qZm6BdLjFFc",
		X5tS256: "b1uYMHKFL1TZeYlC3F6VUH8ApjCf401m_cGCMsxvFlU",
	}

	jwkNext := jwks.Jwk{
		Kid:     "271E7534-C67B-444C-9509-F9A45398EE09",
		Kty:     "RSA",
		Alg:     "RS256",
		Use:     "sig",
		N:       "mrN-PyxusPaAXclD1Og6A4FKnWC_oK2D5IoCURTs6bSsiJezgmAkQSCHcqsIuLJpksxcGyWRE8fhW4VNGxvzFCXBzTzRispm8ExP9GfFNp8gF14ZDfTteQSZqRzrYwhdoaRpDXfJQUnvgOCIVcnJi-tBCcf7TYmeRXIpU4teJPFrKVkb1DLnTPnoUmeDWMietg1PeKluJxVpocOt0vlnpvCtQmM63K1ShnR3a7cqqJnG5lJz7AdeaGfF2Xe_zjdRvYeR1m-cyEt-1-MgjZkwyXglMHOeUdNJV5ekL59_C0YAkmTyp6Sn5QZdT6Z5NrD1jwSGH0zVlDudl8833ui7Hw",
		E:       "AQAB",
		X5c:     []string{"MIIC4TCCAcmgAwIBAgIUTXUqYpCNgcNca+OUOsnIPztEQi0wDQYJKoZIhvcNAQELBQAwADAeFw0yNTA0MTExMzE4MzlaFw0yODAxMDYxMzE4MzlaMAAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCas34/LG6w9oBdyUPU6DoDgUqdYL+grYPkigJRFOzptKyIl7OCYCRBIIdyqwi4smmSzFwbJZETx+FbhU0bG/MUJcHNPNGKymbwTE/0Z8U2nyAXXhkN9O15BJmpHOtjCF2hpGkNd8lBSe+A4IhVycmL60EJx/tNiZ5FcilTi14k8WspWRvUMudM+ehSZ4NYyJ62DU94qW4nFWmhw63S+Wem8K1CYzrcrVKGdHdrtyqomcbmUnPsB15oZ8XZd7/ON1G9h5HWb5zIS37X4yCNmTDJeCUwc55R00lXl6Qvn38LRgCSZPKnpKflBl1Ppnk2sPWPBIYfTNWUO52Xzzfe6LsfAgMBAAGjUzBRMB0GA1UdDgQWBBRX7APTzRcA4uTnGFzl6IxF1uDNizAfBgNVHSMEGDAWgBRX7APTzRcA4uTnGFzl6IxF1uDNizAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQA4cc5/OrC3hUJg9NSdeHcQhBvZefWdAxxOtKQqw3kr7+aYu5YyBb1nVtzsi7wBEAqdhWefZSqjXs6pz4kZa7f78za1+6JUD34TCOKItR3kmSFSc6g3yysgKF9DbeeJz663c8xIwUhkCqJxfjMfLcmTouRbA4qmrQXLaGiQes1dQFLf9ftOs5iRJ0znP+aVl4dfPrsmmpw7pC2M27p8R5dXpfxAzSoRGXOD1zOd7rfRvb50m351AtlRnt5/Qv5ueZU6aXj+4DxcuLwVCCKiENpAoK4rSyt7LamwGS8u8xssjVmwP+XS7+lHz/Pvc6qTofckZd+Dro7oNB/P3rA0t8iH"},
		X5t:     "33eaQDdpm1muofcpGADFWOBOVEo",
		X5tS256: "9dvy1q-_kHs7wok0Yo0X0Trkc2MBnNlTeIoB_yjpvqk",
	}

	tests := []struct {
		description  string
		route        string
		expectedCode int
		expectedJwks server.JwksResponse
	}{
		{
			description:  "Test /certs/default endpoint #1",
			route:        config.GetConfig().ServerConfig.BasePath + "/certs/default",
			expectedCode: 200,
			expectedJwks: server.JwksResponse{
				Keys: []*jwks.Jwk{&jwkNext, &jwkActive, &jwkPrev},
			},
		},
		{
			description:  "Test .../openid-connect/certs endpoint #2",
			route:        "/auth/realms/default/protocol/openid-connect/certs",
			expectedCode: 200,
			expectedJwks: server.JwksResponse{
				Keys: []*jwks.Jwk{&jwkNext, &jwkActive, &jwkPrev},
			},
		},
	}

	jwksConfig := &config.JwksFileConfig{
		UpdateInterval:     0,
		MountedPath:        "./router_testdata/",
		CertFileNameNext:   "next-tls.crt",
		KidFileNameNext:    "next-tls.kid",
		CertFileNameActive: "tls.crt",
		KidFileNameActive:  "tls.kid",
		CertFileNamePrev:   "prev-tls.crt",
		KidFileNamePrev:    "prev-tls.kid",
	}

	srv := server.New()
	jwksProvider, _ := jwks.NewFileProvider(jwksConfig)
	handler := server.NewHandler(jwksProvider)
	srv.RegisterRoutes(handler)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			resp, _ := srv.Test(req, 5)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var jwkSet server.JwksResponse
			err = json.Unmarshal(bodyBytes, &jwkSet)
			if err != nil {
				t.Fatalf("Test failed due to error: %v", err)
			}
			fmt.Printf("XXXXXXXXXXXXX")
			assert.Equalf(t, tt.expectedJwks, jwkSet, "Jwks response does not match expected value")
		})
	}
}

func TestDefaultRealmRoute(t *testing.T) {
	defaultRealm := jwks.DefaultRealm{
		Realm:     "default",
		PublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy0skyMX46fzroq2Ma2pr1iP+Rt+x3IKufm6rf54vwcq/jxYPBajNREM0dtKfjj1p590Gme1+QQW3uS03eK5Rp5CNGGonFrzWsqlYa3dYgHpcZ6UgFGBPJvJCqBnEFP3d7zdg4GKOGDGv+KEM49bKm1qfIvxJ+JpATzv06vNptsGrtygol1rVbWkq8cFZ5mIzSe3Jk0vx8tw3rEint4uG8OHNWqfdHBKblTVjuW2w6cYr7gk6ujm9FswjkZ5us0mgBekw0prLK5bYwNzHERdFtvaCvOIwNZvqwsETQFpQFBwB/7kdEFfuSHbDeG0Mg5/aIikKom2TV+bEy21V6Sw/1QIDAQAB",
	}

	tests := []struct {
		description          string
		route                string
		expectedCode         int
		expectedDefaultRealm jwks.DefaultRealm
	}{
		{
			description:          "Test /issuer/:realm endpoint #1",
			route:                config.GetConfig().ServerConfig.BasePath + "/issuer/default",
			expectedCode:         200,
			expectedDefaultRealm: defaultRealm,
		},
		{
			description:          "Test /auth/realms/:realm endpoint #2",
			route:                "/auth/realms/default",
			expectedCode:         200,
			expectedDefaultRealm: defaultRealm,
		},
	}

	jwksConfig := &config.JwksFileConfig{
		UpdateInterval:     0,
		MountedPath:        "./router_testdata/",
		CertFileNameNext:   "next-tls.crt",
		KidFileNameNext:    "next-tls.kid",
		CertFileNameActive: "tls.crt",
		KidFileNameActive:  "tls.kid",
		CertFileNamePrev:   "prev-tls.crt",
		KidFileNamePrev:    "prev-tls.kid",
	}

	srv := server.New()
	jwksProvider, _ := jwks.NewFileProvider(jwksConfig)
	handler := server.NewHandler(jwksProvider)
	srv.RegisterRoutes(handler)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			resp, _ := srv.Test(req, 5)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.description)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var defaultR jwks.DefaultRealm
			err = json.Unmarshal(bodyBytes, &defaultR)
			if err != nil {
				t.Fatalf("Test failed due to error: %v", err)
			}

			assert.Equalf(t, tt.expectedDefaultRealm, defaultR, "DefaultRealm response does not match expected value")
		})
	}
}
