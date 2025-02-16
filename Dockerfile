# Test stage
FROM golang:1.24.0 AS tester
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go test -v ./...

# Build stage
FROM golang:1.24.0 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go mod tidy
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o eui64-calculator ./cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY eui64-calculator ./
COPY static/ ./static/
COPY ui/ ./ui/
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./eui64-calculator"]