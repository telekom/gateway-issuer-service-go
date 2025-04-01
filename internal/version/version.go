// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package version

import "fmt"

//nolint:gochecknoglobals // this variables should be global to be able to overwrite them via -ldflags
var (
	Version   = "dev"
	BuildDate = "unknown"
)

func GetVersionInfo() string {
	return fmt.Sprintf("Version: %s, Build Date: %s", Version, BuildDate)
}
