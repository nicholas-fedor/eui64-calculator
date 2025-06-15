ARG BASE_IMAGE=golang:1.24.4@sha256:10c131810f80a4802c49cab0961bbe18a16f4bb2fb99ef16deaa23e4246fc817

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
LABEL org.opencontainers.image.title="EUI64 Calculator"
LABEL org.opencontainers.image.description="A tool to calculate EUI-64 network addresses"
LABEL org.opencontainers.image.revision=$VCS_REF
LABEL org.opencontainers.image.source="https://github.com/nicholas-fedor/eui64-calculator"
LABEL org.opencontainers.image.authors="Nicholas Fedor <nick@nickfedor.com>"
LABEL org.opencontainers.image.licenses="GPLv3"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/nickfedor/eui64-calculator"
LABEL org.opencontainers.image.base.name="$BASE_IMAGE"
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@v0.3.833
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o eui64-calculator ./cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian12@sha256:b7b9a6953e7bed6baaf37329331051d7bdc1b99c885f6dbeb72d75b1baad54f9
WORKDIR /app
COPY eui64-calculator .
COPY static/ ./static/
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["./eui64-calculator"]