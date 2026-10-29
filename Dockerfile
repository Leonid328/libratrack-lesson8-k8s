# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/api /app/api
COPY migrations /app/migrations
ENV PORT=8080 MIGRATIONS_PATH=/app/migrations
EXPOSE 8080
USER nobody
ENTRYPOINT ["/app/api"]
