# Weather by Zipcode Service

> This project is part of the [FullCycle](https://fullcycle.com.br/) learning program (Pós-Graduação).

## Table of Contents
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Requirements](#rproject-requirements)
- [License](#license)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) + Docker Compose
- Free API key from [WeatherAPI.com](https://weatherapi.com/)

## Quick Start

Clone the repository, create your local `.env` from the template, then bring up the API

```bash
git clone https://github.com/bianavic/fullcycle_go_auction.git
cd fullcycle_go_auction
cp .env.example .env
docker compose up --build
```

### Environment variables

Configured in `.env` at the project root:

| Variable                   | Example                                      | Description |
|----------------------------|----------------------------------------------|-------------|
| `SERVER_PORT`              | `8080`                                       |             |
| `WEATHER_API_KEY`          | `your-api-key`                               |             |
| `WEATHER_API_UR`           | `https://api.weatherapi.com/v1/current.json` |             |
| `VIACEP_API_URL`           | `https://viacep.com.br/ws`                   |             |

---

## Project Requirements

> See **[docs/CHALLENGE.md](docs/CHALLENGE.md)** for the full assignment (pt-BR).

### Required by the challenge

---

## License

This project was developed for educational purposes as part of the FullCycle program
(Pós-Graduação).