# Weather by Zipcode Service

> This project is part of the [FullCycle](https://fullcycle.com.br/) learning program (Pós-Graduação).

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Local Development](#local-development)
- [API](#api)
- [Project Requirements](#project-requirements)
- [Deployment](#deployment)
- [License](#license)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose
- Free API key from [WeatherAPI.com](https://weatherapi.com/)

## Quick Start

Clone the repository, create your local `.env` from the template, then bring up the API.

```bash
git clone https://github.com/bianavic/fullcycle_go_weather-by-zipcode.git
cd fullcycle_go_weather-by-zipcode
cp .env.example .env
docker compose up --build
```

### Environment variables

Configured in `.env` at the project root:

| Variable          | Example                                      | Description                                                       |
|-------------------|----------------------------------------------|-------------------------------------------------------------------|
| `SERVER_PORT`     | `8080`                                       | Port the HTTP server binds to. Overridden by `PORT` on Cloud Run. |
| `WEATHER_API_KEY` | `your-api-key`                               | Required. WeatherAPI.com API key.                                 |
| `WEATHER_API_URL` | `https://api.weatherapi.com/v1/current.json` | WeatherAPI current-conditions endpoint.                           |
| `VIACEP_API_URL`  | `https://viacep.com.br/ws`                   | viaCEP base URL (CEP -> city lookup).                             |

## Local Development

Without Docker (requires Go 1.26+):

```bash
make run            # go run ./cmd/server
make test           # go test -race -cover ./...
make test-coverage  # writes coverage.out + prints per-func summary
make build          # builds bin/server
make fmt            # gofmt + goimports
```

## API

| Method | Path             | Description                                                       |
|--------|------------------|-------------------------------------------------------------------|
| GET    | `/weather/{cep}` | Returns current temperature in C / F / K for a Brazilian zipcode. |
| GET    | `/health`        | Liveness probe — returns `200 ok`.                                |

### Examples

Success (`200 OK`):

```bash
$ curl -i http://localhost:8080/weather/01001000
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8

{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

Invalid zipcode (`422 Unprocessable Entity`):

```bash
$ curl -i http://localhost:8080/weather/123
HTTP/1.1 422 Unprocessable Entity

invalid zipcode
```

Zipcode not found (`404 Not Found`):

```bash
$ curl -i http://localhost:8080/weather/00000000
HTTP/1.1 404 Not Found

can not find zipcode

```

## Project Requirements

> See **[docs/CHALLENGE.md](docs/CHALLENGE.md)** for the full assignment (pt-BR).

### Required by the challenge

- Receive an 8-digit Brazilian CEP.
- Resolve the city via [viaCEP](https://viacep.com.br/).
- Fetch the current temperature via [WeatherAPI.com](https://weatherapi.com/).
- Return temperatures in Celsius, Fahrenheit, and Kelvin.
- Response contract:
    - `200 OK` -> `{ "temp_C": ..., "temp_F": ..., "temp_K": ... }`
    - `422 Unprocessable Entity` + body `invalid zipcode` when the CEP is not 8 digits.
    - `404 Not Found` + body `can not find zipcode` when the CEP is not registered.
- Deploy to Google Cloud Run.

## Deployment

**Cloud Run URL:** _TBD — populated in the deploy PR._

## License

This project was developed for educational purposes as part of the FullCycle program
(Pós-Graduação).