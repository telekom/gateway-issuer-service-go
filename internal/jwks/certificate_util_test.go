// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package jwks_test

import (
	"crypto/x509"
	"encoding/pem"
	"issuer-service-go/internal/jwks"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testPath = "./certificate_util_testdata"
)

func loadCertificate(t *testing.T, path string) *x509.Certificate {
	certPEM, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read certificate file: %v", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		t.Fatalf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	return cert
}

func TestAlg(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		expectedAlg string
	}{
		{
			description: "verify 'alg' is correctly detected and is RS256",
			certPath:    testPath + "/cert.tls",
			expectedAlg: "RS256",
		},
	}

	for _, test := range tests {
		// when
		alg := jwks.Alg()

		// then
		assert.Equalf(t, test.expectedAlg, alg, test.description)
	}
}

func TestX5c(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		x5cPath     string
	}{
		{
			description: "verify 'x5c' is correctly constructed",
			certPath:    testPath + "/cert.tls",
			x5cPath:     testPath + "/x5c",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)
		x5cVerifier, err := os.ReadFile(test.x5cPath)
		if err != nil {
			t.Fatalf("failed to read x5c verifier file: %v", err)
		}

		// when
		x5c := jwks.X5c(cert)

		// then
		assert.Equalf(t, []string{string(x5cVerifier)}, x5c, test.description)
	}
}

func TestX5t(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		x5t         string
	}{
		{
			description: "verify 'x5t' is correctly constructed",
			certPath:    testPath + "/cert.tls",
			x5t:         "I7VU0PT8AkadM5CUYokqUAnrTuo",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)

		// when
		x5t := jwks.X5t(cert)

		// then
		assert.Equalf(t, test.x5t, x5t, test.description)
	}
}

func TestX5tS256(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		x5tS256     string
	}{
		{
			description: "verify 'x5tS256' is correctly constructed",
			certPath:    testPath + "/cert.tls",
			x5tS256:     "lRWAsHGUm28DvmhZeOriaQ-SymA7NU8plzK5iM2INHo",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)

		// when
		x5tS256 := jwks.X5tS256(cert)

		// then
		assert.Equalf(t, test.x5tS256, x5tS256, test.description)
	}
}

func TestN(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		modulus     string
	}{
		{
			description: "verify 'n' is correctly generated",
			certPath:    testPath + "/cert.tls",
			modulus:     "r9HQyK1Ok474UJX69p4e9dXFXjT7SnJchPXkJ7jYa_6E6_oFCwipcFewWl3nzkKed-tiKmwUM9NXsB9EpeGpfha1HsqesGzzUJow7SbYQge3xW3c6W7S7JrM5Wy3nvZEd4V_YvVd-Q7cvSiQAFjriWsaZFQa3xx2zbSmzowC9ymOXz7ZWQgRf5gszPFzhKbsVRCagWoU5ciwok_iEUZU9DjiUgjdCLY_5HC1771nfxvjF_2zCUN4o7PrqkUbfarud-mgRp-op_n-RsJZisqponLvrmEuCv6yjngp3P8Iz99qgeO4aqywoPbpMAuTp9PACBGifFgPEvz-bdVVTInkHw",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)

		// when
		modulus, _ := jwks.N(cert)

		// then
		assert.Equalf(t, test.modulus, modulus, test.description)
	}
}

func TestE(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		exponent    string
	}{
		{
			description: "verify 'e' is correctly generated",
			certPath:    testPath + "/cert.tls",
			exponent:    "AQAB",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)

		// when
		exponent, _ := jwks.E(cert)

		// then
		assert.Equalf(t, test.exponent, exponent, test.description)
	}
}

func TestPublicKey(t *testing.T) {
	// given
	tests := []struct {
		description string
		certPath    string
		publicKey   string
	}{
		{
			description: "verify 'public_key' is correctly generated",
			certPath:    testPath + "/cert.tls",
			publicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr9HQyK1Ok474UJX69p4e9dXFXjT7SnJchPXkJ7jYa/6E6/oFCwipcFewWl3nzkKed+tiKmwUM9NXsB9EpeGpfha1HsqesGzzUJow7SbYQge3xW3c6W7S7JrM5Wy3nvZEd4V/YvVd+Q7cvSiQAFjriWsaZFQa3xx2zbSmzowC9ymOXz7ZWQgRf5gszPFzhKbsVRCagWoU5ciwok/iEUZU9DjiUgjdCLY/5HC1771nfxvjF/2zCUN4o7PrqkUbfarud+mgRp+op/n+RsJZisqponLvrmEuCv6yjngp3P8Iz99qgeO4aqywoPbpMAuTp9PACBGifFgPEvz+bdVVTInkHwIDAQAB",
		},
	}

	for _, test := range tests {
		cert := loadCertificate(t, test.certPath)

		// when
		exponent, _ := jwks.PublicKey(cert)

		// then
		assert.Equalf(t, test.publicKey, exponent, test.description)
	}
}
