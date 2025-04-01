// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package util

//nolint:gosec // we have to provide 'x5t' in JWK so we are backwards-compatible
import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"
)

// Alg return currently always RS256.
func Alg() string {
	return "RS256"
}

// X5c generates the X.509 certificate chain.
//
// The function takes an X.509 certificate as input and returns a slice of base64-encoded DER certificates.
//
// Parameters:
//
//	cert (*x509.Certificate): The X.509 certificate from which to generate the certificate chain.
//
// Returns:
//
//	[]string: A slice containing the base64-encoded DER representation of the certificate.
func X5c(cert *x509.Certificate) []string {
	return []string{base64.StdEncoding.EncodeToString(cert.Raw)}
}

// X5t generates the SHA-1 thumbprint of the given X.509 certificate.
//
// The function takes an X.509 certificate as input and returns the SHA-1 thumbprint
// of the certificate, encoded as a base64 string.
//
// Parameters:
//
//	cert (*x509.Certificate): The X.509 certificate from which to generate the SHA-1 thumbprint.
//
// Returns:
//
//	string: The SHA-1 thumbprint of the certificate, encoded as a base64 string.
//
//nolint:gosec // we have to provide 'x5t' in JWK so we are backwards-compatible
func X5t(cert *x509.Certificate) string {
	hashSha1 := sha1.Sum(cert.Raw)
	return base64.RawURLEncoding.EncodeToString(hashSha1[:])
}

// X5tS256 generates the SHA-256 thumbprint of the given X.509 certificate.
//
// The function takes an X.509 certificate as input and returns the SHA-256 thumbprint
// of the certificate, encoded as a base64 string.
//
// Parameters:
//
//	cert (*x509.Certificate): The X.509 certificate from which to generate the SHA-256 thumbprint.
//
// Returns:
//
//	string: The SHA-256 thumbprint of the certificate, encoded as a base64 string.
func X5tS256(cert *x509.Certificate) string {
	hashSha256 := sha256.Sum256(cert.Raw)
	return base64.RawURLEncoding.EncodeToString(hashSha256[:])
}

func PublicKey(cert *x509.Certificate) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}

	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	pubKeyStr := string(pubKeyPEM)
	pubKeyStr = strings.ReplaceAll(pubKeyStr, "-----BEGIN PUBLIC KEY-----", "")
	pubKeyStr = strings.ReplaceAll(pubKeyStr, "-----END PUBLIC KEY-----", "")
	pubKeyStr = strings.ReplaceAll(pubKeyStr, "\n", "")
	pubKeyStr = strings.ReplaceAll(pubKeyStr, " ", "")

	return pubKeyStr, nil
}

// N generates the modulus of the RSA public key in base64 URL encoding.
func N(cert *x509.Certificate) (string, error) {
	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("public key is not of type RSA")
	}
	return base64.RawURLEncoding.EncodeToString(rsaPubKey.N.Bytes()), nil
}

// E generates the exponent of the RSA public key in base64 URL encoding.
func E(cert *x509.Certificate) (string, error) {
	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("public key is not of type RSA")
	}

	// Convert the exponent (E) to a byte slice
	eBytes := new(big.Int).SetInt64(int64(rsaPubKey.E)).Bytes()

	// Encode the exponent in Base64 URL encoding
	return base64.RawURLEncoding.EncodeToString(eBytes), nil
}
