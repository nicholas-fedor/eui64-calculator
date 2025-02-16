# EUI64 Calculator with Traefik

This runs the application with Traefik in front of it as a reverse proxy providing TLS encryption via Let's Encrypt SSL Certificates.

## Usage

The example Docker Compose stack is specific to domains using Cloudflare's DNS services.

Feel free to reference Traefik's [documentation](https://doc.traefik.io/traefik/) to change this for your particular use case.

The `traefik.yaml` file is set by default to obtain a Let's Encrypt staging certificate. Feel free to disable this and enable the production certificate request if everything is working.

## Installation

1. Create the following directory and file structure. Copy-paste from the repository, as needed:

    ```console
    Docker
    ├─── docker-compose.yaml
    ├─── .env
    └───Traefik
        ├─── traefik.yaml
        ├───Certs
        |   └─── acme.json
        ├───Configs
        └───Secrets
            └─── CLOUDFLARE_DNS_API_TOKEN
    ```

2. Set permissions on `acme.json` certificates file:

    ```console
    chmod 600 ./Traefik/Certs/acme.json
    ```

3. Go to Cloudflare's [website](https://dash.cloudflare.com/profile/api-tokens) to generate an API token with the `Edit zone DNS` template.

4. Update the `.env` file with the domain name that you're going to be using.

    > Don't forget to add the appropriate DNS record entries either in your `hosts` file or DNS resolver.

5. Pull the container images, which also helps to check your `docker-compose.yaml` file is accurate:

    ```console
    docker compose pull
    ```

6. If everything looks good, run the `docker-compose.yaml` stack:

    ```console
    docker compose up -d
    ```

7. Access the application at the URL that you specified, i.e. <https://eui64-calculator.example.com>
