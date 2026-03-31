<!-- markdownlint-disable -->
<div align="center">

# EUI-64 Calculator

A EUI-64 address calculator implemented in Go, HTMX, and Templ.

Inspired by [ThePrincelle's EUI64-Calculator](https://github.com/ThePrincelle/EUI64-Calculator)

### Also available at <https://eui64-calculator.nickfedor.com>

![EUI-64 Calculator Screenshot](./.github/assets/eui64-calculator_screenshot.png)
<br/><br/>
<!-- markdownlint-restore -->

  [![CircleCI](https://dl.circleci.com/status-badge/img/gh/nicholas-fedor/eui64-calculator/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/nicholas-fedor/eui64-calculator/tree/main)
  [![codecov](https://codecov.io/gh/nicholas-fedor/eui64-calculator/branch/main/graph/badge.svg)](https://codecov.io/gh/nicholas-fedor/eui64-calculator)
  [![GoDoc](https://godoc.org/github.com/nicholas-fedor/eui64-calculator?status.svg)](https://godoc.org/github.com/nicholas-fedor/eui64-calculator)
  [![Go Report Card](https://goreportcard.com/badge/github.com/nicholas-fedor/eui64-calculator)](https://goreportcard.com/report/github.com/nicholas-fedor/eui64-calculator)
  [![latest version](https://img.shields.io/github/tag/nicholas-fedor/eui64-calculator.svg)](https://github.com/nicholas-fedor/eui64-calculator/releases)
  [![GPLv3 License](https://img.shields.io/github/license/nicholas-fedor/eui64-calculator.svg)](https://www.gnu.org/licenses/gpl-3.0)
  [![Codacy Badge](https://app.codacy.com/project/badge/Grade/1c48cfb7646d4009aa8c6f71287670b8)](https://www.codacy.com/gh/nicholas-fedor/eui64-calculator/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=nicholas-fedor/eui64-calculator&amp;utm_campaign=Badge_Grade)
  [![All Contributors](https://img.shields.io/github/all-contributors/nicholas-fedor/eui64-calculator)](#contributors)
  [![Pulls from DockerHub](https://img.shields.io/docker/pulls/nickfedor/eui64-calculator.svg)](https://hub.docker.com/r/nickfedor/eui64-calculator)

</div>

## Overview

This project provides a simple tool for calculating an EUI-64 IPv6 address using a MAC addresses and IPv6 Prefix.

## Usage

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

- Running the [Basic Template](/examples/docker-compose.yaml):

    ```console
    docker compose -f ./Docker/compose.yaml up -d
    ```

- Traefik Reverse Proxy [example](/examples/Traefik/README.md)

### Running Locally

#### Prerequisites

- Go: <https://go.dev/doc/install>
- Templ: `go install github.com/a-h/templ/cmd/templ@latest`

#### Installation

1. Clone the repository:

    ```console
    git clone https://github.com/nicholas-fedor/eui64-calculator.git
    ```

2. Enter the repository:

    ```console
    cd eui64-calculator
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
├── .github
│   ├── workflows
│   │   ├── create-manifests.yaml
│   │   ├── lint-go.yaml
│   │   ├── build.yaml
│   │   ├── clean-cache.yaml
│   │   ├── pull-request.yaml
│   │   ├── release.yaml
│   │   ├── security.yaml
│   │   └── test.yaml
│   ├── ISSUE_TEMPLATE
│   │   ├── bug_report.yaml
│   │   ├── config.yaml
│   │   └── feature_request.yaml
│   ├── assets
│   │   ├── eui64-calculator_screenshot.png
│   │   └── eui64-calculator_social-preview_1280x640.png
│   └── renovate.json
├── build
│   ├── docker
│   │   ├── Dockerfile
│   │   └── .dockerignore
│   ├── golangci-lint
│   │   └── golangci.yaml
│   └── goreleaser
│       └── goreleaser.yaml
├── cmd
│   └── server
│       ├── static
│       │   ├── favicon.ico
│       │   └── styles.css
│       ├── main.go
│       └── main_test.go
├── internal
│   ├── ui
│   │   ├── doc.go
│   │   ├── generate.go
│   │   ├── home.templ
│   │   ├── layout.templ
│   │   ├── result.templ
│   │   └── ui_test.go
│   ├── eui64
│   │   ├── eui64.go
│   │   └── eui64_test.go
│   ├── handlers
│   │   ├── handlers.go
│   │   └── handlers_test.go
│   └── validators
│       ├── ipv6_prefix_validator.go
│       ├── ipv6_prefix_validator_test.go
│       ├── mac_validator.go
│       └── mac_validator_test.go
├── .all-contributorsrc
├── .circleci
│   └── config.yml
├── .codacy.yml
├── .gitattributes
├── .gitignore
├── LICENSE
├── README.md
├── examples
│   ├── Traefik
│   │   ├── .env
│   │   ├── README.md
│   │   ├── docker-compose.yaml
│   │   └── traefik.yaml
│   └── docker-compose.yaml
├── go.mod
└── go.sum
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
    docker build -f docker/Dockerfile-dev -t eui64-calculator-dev .
    ```

- Running the image locally:

    ```console
    docker run -it -p 8080:8080 eui64-calculator-dev
    ```

### Notes

- The Dockerfile uses `gcr.io/distroless/static-debian12` as the final runtime image for the application. This results in a minimal container image without a shell or other features typical of other container images.

- I opted to hardcode Gin's release mode to avoid redundant environment variables. This can be easily commented out in the `cmd/server/main.go` file.

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->

<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/nicholas-fedor"><img src="https://avatars2.githubusercontent.com/u/71477161?v=4?s=100" width="100px;" alt="Nicholas Fedor"/><br /><sub><b>Nicholas Fedor</b></sub></a><br /><a href="https://github.com/nicholas-fedor/eui64-calculator/commits?author=nicholas-fedor" title="Code">💻</a> <a href="https://github.com/nicholas-fedor/eui64-calculator/commits?author=nicholas-fedor" title="Documentation">📖</a> <a href="#maintenance-nicholas-fedor" title="Maintenance">🚧</a> <a href="https://github.com/nicholas-fedor/eui64-calculator/pulls?q=is%3Apr+reviewed-by%3Anicholas-fedor" title="Reviewed Pull Requests">👀</a></td>
    </tr>
  </tbody>
</table>

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
