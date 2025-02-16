# Test stage
FROM golang:1.24.0 AS tester
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go test -v ./...

# Build stage
FROM golang:1.24.0 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o eui64-calculator ./cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian12:latest
WORKDIR /app
COPY --from=builder /app/eui64-calculator .
COPY --from=builder /app/static/ ./static/
COPY --from=builder /app/ui/ ./ui/
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./eui64-calculator"]