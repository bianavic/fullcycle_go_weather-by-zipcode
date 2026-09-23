# Weather by Zipcode Service

> This project is part of the [FullCycle](https://fullcycle.com.br/) learning program (Pós-Graduação).

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start & Local Testing (Docker)](#quick-start--local-testing-docker)
- [API Endpoints & Validation Scenarios](#api-endpoints--validation-scenarios)
- [Project Requirements](#project-requirements)
- [Cloud Run Deployment & Evaluator Guide](#cloud-run-deployment--evaluator-guide)
- [License](#license)

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose (optional, for local container testing).
- Free API key from [WeatherAPI.com](https://www.weatherapi.com/).

---

## Quick Start & Local Testing (Docker)

To run the application locally using Docker Compose, clone the repository, set up your `.env` file with a valid WeatherAPI key, and start the container:

```bash
git clone https://github.com/bianavic/fullcycle_go_weather-by-zipcode.git
cd fullcycle_go_weather-by-zipcode

# Create your local environment file
cp .env.example .env
# Edit .env and fill in your WEATHER_API_KEY=your_key_here

# Build and start the container
make docker-up         # or: docker compose up --build
```

To stop the local container:

```bash
make docker-down       # or: docker compose down
```

---

## API Endpoints & Validation Scenarios

The service runs on port `8080` and exposes the following endpoints to fulfill the challenge criteria:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Root endpoint — confirms the API is running and gives usage instructions. `Weather API is running! Use /weather/{cep}`. |
| `GET` | `/health` | Liveness probe — returns `200 ok`. |
| `GET` | `/weather/{cep}` | Receives an 8-digit CEP, queries location via viaCEP, fetches weather, and returns temperatures in Celsius, Fahrenheit, and Kelvin. |

### Testing via `curl` (Local Example)

#### 0. Root Endpoint (`200 OK`)

```bash
curl -i http://localhost:8080/
```

#### 1. Success (`200 OK`)

```bash
curl -i http://localhost:8080/weather/01001000
```

Expected response (`200`):

```json
{"temp_C":22.4,"temp_F":72.32,"temp_K":295.4}
```

#### 2. Invalid Zipcode Format (`422 Unprocessable Entity`)

Triggered when the CEP does not contain exactly 8 digits.

```bash
curl -i http://localhost:8080/weather/123
```

Expected response (`422`):

```text
invalid zipcode
```

#### 3. Zipcode Not Found (`404 Not Found`)

Triggered when the 8-digit CEP format is valid, but the location does not exist in viaCEP.

```bash
curl -i http://localhost:8080/weather/00000000
```

Expected response (`404`):

```text
can not find zipcode
```

#### 4. Health Check (`200 OK`)

```bash
curl -i http://localhost:8080/health
```

Expected response (`200`):

```text
ok
```

---

## Project Requirements

Derived from [CHALLENGE.md](docs/CHALLENGE.md):

- Receive an 8-digit Brazilian CEP.
- Resolve city/location using the [viaCEP API](https://viacep.com.br/).
- Fetch current weather using [WeatherAPI.com](https://www.weatherapi.com/).
- Return temperatures formatted in **Celsius**, **Fahrenheit**, and **Kelvin**.
- Automated tests covering system logic (`make test`).
- Deployed to Google Cloud Run (free tier) — see [Cloud Run Deployment & Evaluator Guide](#cloud-run-deployment--evaluator-guide).

---

## Cloud Run Deployment & Evaluator Guide

> **Note for Evaluators:** The service is fully deployed on Google Cloud Run (free tier) and is publicly accessible without requiring local setup, API keys, or manual container building. You can test it directly using the live URL below.

These four scenarios map directly to the challenge's response-format requirements.

| Scenario                     | Expected status | Expected body                              |
|-------------------------------|:----------------:|----------------------------------------------|
| Valid CEP                     | `200`            | `{"temp_C":...,"temp_F":...,"temp_K":...}`    |
| Malformed CEP (not 8 digits)  | `422`            | `invalid zipcode`                             |
| Well-formed but unknown CEP   | `404`            | `can not find zipcode`                        |
| Health check                  | `200`            | `ok`                                          |

### Live Cloud Run Service URL

**Base URL:** `https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app`

### How to Validate the Deployed Service

You can copy and paste these commands directly into your terminal to test all required evaluation criteria against the live environment:

**1. Test Root Endpoint:**

```bash
curl -i "https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app/"
```

**2. Test liveness / health check:**

```bash
curl -i "https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app/health"
```

**3. Test success scenario (`200 OK` with temperatures):**
*(You can also open this link directly in your browser.)*

```bash
curl -i "https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app/weather/01001000"
```

**4. Test invalid zipcode scenario (`422 Unprocessable Entity`):**

```bash
curl -i "https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app/weather/123"
```

**5. Test zipcode not found scenario (`404 Not Found`):**

```bash
curl -i "https://weather-by-zipcode-kxqudjaa3a-uc.a.run.app/weather/00000000"
```

---

## License

This project was developed for educational purposes as part of the FullCycle program (Pós-Graduação).
