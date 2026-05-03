ARG GO_VERSION=1.26.2
ARG NODE_VERSION=22-alpine
ARG ALPINE_VERSION=3.22

FROM node:${NODE_VERSION} AS web-builder
WORKDIR /src/web

COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

FROM golang:${GO_VERSION}-alpine AS go-builder
WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /src/web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/chattery \
    ./cmd/main.go

FROM alpine:${ALPINE_VERSION} AS runtime

RUN addgroup -S chattery && adduser -S -D -H -G chattery chattery

WORKDIR /app

COPY --from=go-builder /out/chattery /app/chattery

ENV APP_DEBUG=false \
    HTTP_HOST=0.0.0.0 \
    HTTP_PORT=8080

USER chattery
EXPOSE 8080

ENTRYPOINT ["/app/chattery"]
