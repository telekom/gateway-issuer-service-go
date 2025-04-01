// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package jwks

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"issuer-service-go/internal/config"
	"issuer-service-go/internal/util"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	schedulerName = "JWKS File Provider - Scheduler"
)

type Jwk struct {
	Kid       string   `json:"kid"`
	Kty       string   `json:"kty"`
	Alg       string   `json:"alg"`
	Use       string   `json:"use"`
	N         string   `json:"n"`
	E         string   `json:"e"`
	X5c       []string `json:"x5c"`
	X5t       string   `json:"x5t"`
	X5tS256   string   `json:"x5t#S256"`
	PublicKey string   `json:"-"`
}

type DefaultRealm struct {
	Realm     string `json:"realm"`
	PublicKey string `json:"public_key"`
}

type FileProvider struct {
	config *config.JwksFileConfig

	certsCacheMap map[config.Type]*Jwk

	cacheMutex *sync.Mutex

	isSchedulerRunning bool
}

func NewFileProvider(jwksConfig *config.JwksFileConfig) (*FileProvider, error) {
	fp := &FileProvider{
		config:        jwksConfig,
		certsCacheMap: make(map[config.Type]*Jwk),
		cacheMutex:    &sync.Mutex{},
	}
	if err := initialize(fp); err != nil {
		return nil, fmt.Errorf("failed to initialize FileProvider: %w", err)
	}
	return fp, nil
}

func (fp *FileProvider) GetJwks() []*Jwk {
	fp.cacheMutex.Lock()
	defer fp.cacheMutex.Unlock()

	keyOrder := []config.Type{config.Next, config.Active, config.Previous}

	values := make([]*Jwk, 0, len(keyOrder))
	for _, key := range keyOrder {
		if jwk, exists := fp.certsCacheMap[key]; exists {
			values = append(values, jwk)
		}
	}
	return values
}

func (fp *FileProvider) GetDefaultRealm(realm string) *DefaultRealm {
	fp.cacheMutex.Lock()
	defer fp.cacheMutex.Unlock()

	defaultRealm := &DefaultRealm{
		Realm:     realm,
		PublicKey: fp.certsCacheMap[config.Active].PublicKey,
	}

	return defaultRealm
}

func (fp *FileProvider) IsSchedulerRunning() bool {
	return fp.isSchedulerRunning
}

func initialize(fp *FileProvider) error {
	log.Info().Msgf("initializing JWKS cache...")
	fp.cacheMutex.Lock()
	defer fp.cacheMutex.Unlock()

	if err := updateCerts(fp); err != nil {
		return fmt.Errorf("failed to update certificates: %w", err)
	}

	log.Info().Msgf("JWKS cache is initialized")
	startScheduler(fp)
	return nil
}

func updateCerts(fp *FileProvider) error {
	jwkNext, err := generateCertInfo(fp.config, config.Next)
	if err != nil {
		return err
	}
	jwkActive, err := generateCertInfo(fp.config, config.Active)
	if err != nil {
		return err
	}
	jwkPrev, err := generateCertInfo(fp.config, config.Previous)
	if err != nil {
		return err
	}

	fp.certsCacheMap = make(map[config.Type]*Jwk)

	addJwkToCache(fp, config.Active, jwkActive)
	addJwkToCache(fp, config.Previous, jwkPrev)
	addJwkToCache(fp, config.Next, jwkNext)

	return nil
}

func addJwkToCache(fp *FileProvider, certType config.Type, jwk *Jwk) {
	var found bool
	for _, value := range fp.certsCacheMap {
		if value.Kid == jwk.Kid {
			found = true
			break
		}
	}

	if !found {
		fp.certsCacheMap[certType] = jwk
	} else {
		log.Debug().Msgf("JWK with kid %s already exists in cache", jwk.Kid)
	}
}

func generateCertInfo(config *config.JwksFileConfig, certType config.Type) (*Jwk, error) {
	certFile := config.GetCertFile(certType)
	certByteArray, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	kidFile := config.GetKidFile(certType)
	kidByteArray, err := os.ReadFile(kidFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(certByteArray)
	if block == nil {
		return nil, errors.New("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	exponent, err := util.E(cert)
	if err != nil {
		return nil, fmt.Errorf("unable to create JWK: %w", err)
	}

	modulus, err := util.N(cert)
	if err != nil {
		return nil, fmt.Errorf("unable to create JWK: %w", err)
	}

	// Extract and format the public key
	publicKeyString, err := util.PublicKey(cert)
	if err != nil {
		return nil, fmt.Errorf("unable to read Public Key: %w", err)
	}

	jwk := Jwk{
		Kid:       string(kidByteArray),
		Kty:       "RSA",
		Alg:       util.Alg(),
		Use:       "sig",
		E:         exponent,
		N:         modulus,
		X5c:       util.X5c(cert),
		X5t:       util.X5t(cert),
		X5tS256:   util.X5tS256(cert),
		PublicKey: publicKeyString,
	}

	return &jwk, nil
}

func startScheduler(fp *FileProvider) {
	if fp.config.UpdateInterval == 0 {
		log.Info().Msgf("%s is deactivated", schedulerName)
		return
	}

	log.Info().Msgf("starting %s ...", schedulerName)
	ticker := time.NewTicker(time.Duration(fp.config.UpdateInterval) * time.Second)

	go func() {
		for range ticker.C {
			executeTask(fp)
		}
	}()

	fp.isSchedulerRunning = true
	log.Info().Msgf("%s started", schedulerName)
}

func executeTask(fp *FileProvider) {
	log.Info().Msg("updating the certificates from mounted files...")
	err := updateCerts(fp)
	if err != nil {
		log.Error().Msgf("failed to update certificate: %v", err)
	}
	log.Info().Msg("certificates were updated successfully")
}
