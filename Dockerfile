# SPDX-FileCopyrightText: 2025 Deutsche Telekom AG
#
# SPDX-License-Identifier: Apache-2.0

FROM alpine:latest

WORKDIR /app

RUN addgroup -g 1000 -S app
RUN adduser -u 1000 -D -H -S -G app app

COPY --chown=app:app issuer-service /app/issuer-service

USER 1000:1000

EXPOSE 8080

CMD ["./issuer-service"]
