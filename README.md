<div align="center">

# EUI-64 Calculator

A EUI-64 address calculator implemented in Go, HTMX, and Templ.

This project was inspired by [ThePrincelle's EUI64-Calculator](https://github.com/ThePrincelle/EUI64-Calculator)
<br/><br/>

  [![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/EUI64-Calculator/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/EUI64-Calculator/tree/main)
  [![codecov](https://codecov.io/gh/nicholas-fedor/EUI64-Calculator/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/EUI64-Calculator)
  [![GoDoc](https://godoc.org/github.com/nicholas-fedor/EUI64-Calculator?status.svg)](https://godoc.org/github.com/nicholas-fedor/EUI64-Calculator)
  [![Go Report Card](https://goreportcard.com/badge/github.com/nicholas-fedor/EUI64-Calculator)](https://goreportcard.com/report/github.com/nicholas-fedor/EUI64-Calculator)
  [![latest version](https://img.shields.io/github/tag/nicholas-fedor/EUI64-Calculator.svg)](https://github.com/nicholas-fedor/EUI64-Calculator/releases)
  [![GPLv3 License](https://img.shields.io/github/license/nicholas-fedor/EUI64-Calculator.svg)](https://www.gnu.org/licenses/gpl-3.0)
  [![Codacy Badge](https://app.codacy.com/project/badge/Grade/1c48cfb7646d4009aa8c6f71287670b8)](https://www.codacy.com/gh/nicholas-fedor/EUI64-Calculator/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/EUI64-Calculator&amp;utm_campaign=Badge_Grade)
  [![All Contributors](https://img.shields.io/github/all-contributors/nicholas-fedor/EUI64-Calculator)](#contributors)
  [![Pulls from DockerHub](https://img.shields.io/docker/pulls/nickfedor/EUI64-Calculator.svg)](https://hub.docker.com/r/nickfedor/EUI64-Calculator)
</div>

## Overview

This project provides a simple tool for calculating an EUI-64 IPv6 address using a MAC addresses and IPv6 Prefix.

### Features

- **EUI-64 Calculation**: Convert a 48-bit MAC address into a 64-bit EUI-64 format.
- **IPv6 Address Generation**: Combine the EUI-64 with a user-provided IPv6 prefix.
- **Web Interface**: User-friendly interface for input and result display using HTMX for dynamic content loading.
- **Docker Support**: Containerized deployment for easy setup and scalability.

### Usage

1. Enter a MAC Address in the format `xx-xx-xx-xx-xx-xx`.
2. Enter an IPv6 Prefix.
3. Click `Calculate` to see the results.

## Getting Started

### Docker Deployment

#### Quick Start

```console
docker run -d --name eui64-calculator nickfedor/eui64-calculator:latest
```

#### Docker Compose

- Running the [Basic Template](/docker/docker-compose.yaml):

    ```console
    docker compose -f ./Docker/compose.yaml up -d
    ```

- Traefik Reverse Proxy [example](/docker/Examples/Traefik/README.md)

### Running Locally

#### Prerequisites

- Go: <https://go.dev/doc/install>
- Templ: `go install github.com/a-h/templ/cmd/templ@latest`

#### Installation

1. Clone the repository:

    ```console
    git clone https://github.com/nicholas-fedor/EUI64-Calculator.git
    ```

2. Enter the repository:

    ```console
    cd EUI64-Calculator
    ```

3. Install Dependencies:

    ```console
    go mod download
    ```

4. Generate Templates:

    ```console
    templ generate
    ```

5. Run the Server:

    ```console
    go run ./cmd/server/main.go
    ```

6. The application will be accessible at <http://localhost:8080/>

## Development

### Project Structure

```console
.
├── cmd
│   └── server
│       ├── main.go
│       └── main_test.go
├── docker
│   ├── docker-compose.yaml
│   ├── Dockerfile
│   └── Examples
│       └── Traefik
├── internal
│   ├── eui64
│   │   ├── eui64.go
│   │   └── eui64_test.go
│   └── handlers
│       ├── handlers.go
│       └── handlers_test.go
├── static
│   ├── favicon.ico
│   └── styles.css
└── ui
    ├── home.templ
    ├── home_templ.go
    ├── layout.templ
    ├── layout_templ.go
    ├── result.templ
    ├── result_templ.go
    └── ui_test.go
```

### Dependencies

- Golang: <https://go.dev/doc>
- gin-gonic/gin: <https://github.com/gin-gonic/gin>
- Templ: <https://github.com/a-h/templ>
- HTMX: <https://htmx.org/docs>

### IDE Support

If you're using VSCode, I've included an `extensions.json` file with recommended extensions.

### Managing Templ files

- Installing the Templ CLI

    ```console
    go install github.com/a-h/templ/cmd/templ@latest
    ```

- Rebuilding `.templ.go` files after updates to `.templ` files (run from the project's root directory)

    Linux:

    ```console
    rm ./ui/*_templ.go && templ generate
    ```

    Windows:

    ```console
    del ui\*_templ.go && templ generate
    ```

### Testing

- Unit Tests:

    ```console
    go test ./...
    ```

- Docker Test Stage:

    The Dockerfile includes a test stage to ensure all tests pass before building the production image.

### Docker

- Rebuilding the Docker image:

    ```console
    docker build -f ./Docker/Dockerfile -t eui64-calculator:latest .
    ```

- Running the image locally:

    ```console
    docker run -p 8080:8080 eui64-calculator:latest
    ```

### Notes

- The Dockerfile uses `gcr.io/distroless/static-debian12` as the final runtime image for the application. This results in a minimal container image without a shell or other features typical of other container images.

- I opted to hardcode Gin's release mode to avoid redundant environment variables. This can be easily commented out in the `cmd/server/main.go` file.

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Contributing

This was a weekend project and there's plenty of opportunity for improvement.

If you feel like contributing, please:

- Fork the repo
- Create your feature branch: `git checkout -b feature/AmazingFeature`
- Commit your changes: `git commit -m "Add some AmazingFeature"`
- Push to the branch: `git push origin feature/AmazingFeature`
- Open a pull request

## License

This project is licensed under the GNU GPLv3 license - see the [LICENSE](#license) file for details.
