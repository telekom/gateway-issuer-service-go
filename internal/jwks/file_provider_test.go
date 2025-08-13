// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package jwks_test

import (
	"issuer-service-go/internal/config"
	"issuer-service-go/internal/jwks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFileProvider(t *testing.T) {
	tests := []struct {
		name   string
		config *config.JwksFileConfig
		err    bool
	}{
		{
			name: "no error with valid config",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: false,
		},
		{
			name: "error with invalid CertFileNameNext",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "invalid-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid KidFileNameNext",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "invalid-tls.kid",
				CertFileNameActive: "tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid CertFileNameActive",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "invalid-tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid KidFileNameActive",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "next-tls.crt",
				KidFileNameActive:  "invalid-tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid CertFileNamePrev",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "next-tls.crt",
				KidFileNameActive:  "next-tls.kid",
				CertFileNamePrev:   "invalid-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid KidFileNamePrev",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "next-tls.crt",
				KidFileNameActive:  "next-tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "invalid-tls.kid",
			},
			err: true,
		},
		{
			name: "error with invalid MountedPath",
			config: &config.JwksFileConfig{
				UpdateInterval:     0,
				MountedPath:        "./file_provider_invalid",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "next-tls.crt",
				KidFileNameActive:  "next-tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwks.NewFileProvider(tt.config)
			assert.Equalf(t, tt.err, err != nil, "expected error: %v, got: %v", tt.err, err)
		})
	}
}

func TestGetJwks(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.JwksFileConfig
		err      bool
		certsCnt int
	}{
		{
			name: "validate that the scheduler is running",
			config: &config.JwksFileConfig{
				UpdateInterval:     1,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls.kid",
				CertFileNameActive: "tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err:      false,
			certsCnt: 3,
		},
		{
			name: "validate that the same kid is not stored twice",
			config: &config.JwksFileConfig{
				UpdateInterval:     1,
				MountedPath:        "./file_provider_testdata",
				CertFileNameNext:   "next-tls.crt",
				KidFileNameNext:    "next-tls-samekid.kid",
				CertFileNameActive: "tls.crt",
				KidFileNameActive:  "tls.kid",
				CertFileNamePrev:   "prev-tls.crt",
				KidFileNamePrev:    "prev-tls.kid",
			},
			err:      false,
			certsCnt: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwksProvider, _ := jwks.NewFileProvider(tt.config)
			time.Sleep(3 * time.Second)
			assert.Truef(t, jwksProvider.IsSchedulerRunning(), "expected scheduler to be running, but it is not")

			jwKeySet := jwksProvider.GetJwks()
			assert.Lenf(
				t,
				jwKeySet,
				tt.certsCnt,
				"expected JWKS to has length %d, but got %d",
				tt.certsCnt,
				len(jwKeySet),
			)
		})
	}
}
