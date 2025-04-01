// SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
//
// SPDX-License-Identifier: Apache-2.0

package jwks

type Provider interface {
	GetJwks() []*Jwk
	GetDefaultRealm(realm string) *DefaultRealm
}
