ARG BASE_IMAGE=golang:1.24.0

# Test stage
FROM $BASE_IMAGE AS tester
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY /cmd /cmd
COPY /internal /internal
COPY /static /static
COPY /ui /ui
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN templ generate
RUN go test -v ./...

# Build stage
FROM $BASE_IMAGE AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o eui64-calculator ./cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian12@sha256:3f2b64ef97bd285e36132c684e6b2ae8f2723293d09aae046196cca64251acac
WORKDIR /app
COPY eui64-calculator .
COPY static/ ./static/
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./eui64-calculator"]