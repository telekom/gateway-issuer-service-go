// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals // ignore this linter error
var current *Config

func GetConfig() *Config {
	if current == nil {
		initializeConfig()
	}

	return current
}

func initializeConfig() {
	current = &Config{}
	if err := env.Parse(current); err != nil {
		log.Fatal().Msgf("failed to parse config: %v", err)
	}

	log.Info().Msgf("config parsed from env: %+v", current)

	level, err := zerolog.ParseLevel(strings.ToLower(current.LogLevel))
	if err != nil {
		log.Warn().Msgf("invalid log level '%s', defaulting to 'info'", current.LogLevel)
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Info().Msgf("log level set to '%s'", level.String())
}
