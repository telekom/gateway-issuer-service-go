// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"path"
	"time"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL,expand" envDefault:"info"` // Log level of the application

	//nolint:golines // ignore this linter error
	GracefulShutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT,expand" envDefault:"5s"`     // Timeout in seconds for graceful shutdown
	IssuerURL               string        `env:"ISSUER_URL,expand" envDefault:"http://localhost:8080"` // URL of the issuer that should be used in the JWKS endpoint

	ServerConfig ServerConfig
	JwksConfig   JwksFileConfig
}

//nolint:golines // ignore this linter error
type ServerConfig struct {
	Port     int    `env:"SERVER_PORT,expand" envDefault:"8081"`      // Port the server should listen on
	BasePath string `env:"API_BASE_PATH,expand" envDefault:"/api/v1"` // Base path of the API
}

//nolint:golines // ignore this linter error
type JwksFileConfig struct {
	UpdateInterval     int    `env:"CERT_UPDATE_INTERVAL,expand" envDefault:"10"`     // Interval in seconds in which the certificates should be updated. If 0 scheduler is deactivated at all
	MountedPath        string `env:"CERT_MOUNT_PATH,expand,required"`                 // Path to the directory where the certificates are mounted
	CertFileNameNext   string `env:"CERT_FILE_NEXT,expand" envDefault:"next-tls.crt"` // Name of the certificate file that should be used in the next rotation
	KidFileNameNext    string `env:"KID_FILE_NEXT,expand" envDefault:"next-tls.kid"`  // Name of the key ID file that should be used in the next rotation
	CertFileNameActive string `env:"CERT_FILE_ACTIVE,expand" envDefault:"tls.crt"`    // Name of the certificate file that should be used currently
	KidFileNameActive  string `env:"KID_FILE_ACTIVE,expand" envDefault:"tls.kid"`     // Name of the key ID file that should be used currently
	CertFileNamePrev   string `env:"CERT_FILE_PREV,expand" envDefault:"prev-tls.crt"` // Name of the certificate file that should be used to verify the signature of JWTs that were signed with a key that is not the current one
	KidFileNamePrev    string `env:"KID_FILE_PREV,expand" envDefault:"prev-tls.kid"`  // Name of the key ID file that should be used to verify the signature of JWTs that were signed with a key that is not the current one
}

type Type int

const (
	Next Type = iota
	Active
	Previous
)

func (c *JwksFileConfig) GetCertFile(jwksType Type) string {
	switch jwksType {
	case Next:
		return path.Join(c.MountedPath, c.CertFileNameNext)
	case Active:
		return path.Join(c.MountedPath, c.CertFileNameActive)
	case Previous:
		return path.Join(c.MountedPath, c.CertFileNamePrev)
	}
	return ""
}

func (c *JwksFileConfig) GetKidFile(jwksType Type) string {
	switch jwksType {
	case Next:
		return path.Join(c.MountedPath, c.KidFileNameNext)
	case Active:
		return path.Join(c.MountedPath, c.KidFileNameActive)
	case Previous:
		return path.Join(c.MountedPath, c.KidFileNamePrev)
	}
	return ""
}
