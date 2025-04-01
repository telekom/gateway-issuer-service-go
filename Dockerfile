# SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VERSION=1.23.8
FROM golang:${GO_VERSION}-alpine

WORKDIR /app

COPY . .

RUN go build -ldflags="-X 'internal/version.Version=${VERSION}' -X 'internal/version.BuildDate=$(date -u +'%Y-%m-%dT%H:%M:%SZ')'" -o issuer-service cmd/api/main.go

RUN addgroup -g 1000 -S app
RUN adduser -u 1000 -D -H -S -G app app

USER 1000:1000

EXPOSE 8080

CMD ["./issuer-service"]
