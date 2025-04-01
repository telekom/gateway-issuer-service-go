// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/caarlos0/env/v11"
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
}
