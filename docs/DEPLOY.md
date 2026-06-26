# Deployment — Google Cloud Run

The service is deployed as a container image hosted in **Google Artifact Registry** and
served by **Cloud Run** (fully managed). The image is built from the project `Dockerfile`
(multi-stage `golang:1.26` → `distroless/static:nonroot`).

## Prerequisites

Install and authenticate the Google Cloud CLI:

```bash
gcloud auth login
gcloud config set project <PROJECT_ID>
```

Enable the required APIs (one-time per project):

```bash
gcloud services enable \
  run.googleapis.com \
  artifactregistry.googleapis.com
```

Confirm Docker is installed locally — the deploy script uses `docker build` / `docker push`.

## One-time setup

### 1. Create the Artifact Registry repository

```bash
gcloud artifacts repositories create weather-by-zipcode \
  --repository-format=docker \
  --location=us-central1 \
  --description="Weather by Zipcode container images"
```

### 2. Configure Docker to push to Artifact Registry

```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

## Deploy

### Option A — script (recommended)

```bash
export PROJECT_ID=<your-project-id>
export WEATHER_API_KEY=<your-weatherapi-key>
./scripts/deploy.sh
```

The script builds the image (tagged with the current git short SHA), pushes it to
Artifact Registry, deploys the revision to Cloud Run with public access, and prints
the service URL.

Overridable env vars (with defaults):

| Variable          | Default                                        |
|-------------------|------------------------------------------------|
| `REGION`          | `us-central1`                                  |
| `SERVICE`         | `weather-by-zipcode`                           |
| `REPO`            | `weather-by-zipcode`                           |
| `IMAGE_TAG`       | git short SHA (or `latest`)                    |
| `WEATHER_API_URL` | `https://api.weatherapi.com/v1/current.json`   |
| `VIACEP_API_URL`  | `https://viacep.com.br/ws`                     |

### Option B — manual

```bash
PROJECT_ID=<your-project-id>
REGION=us-central1
SERVICE=weather-by-zipcode
REPO=weather-by-zipcode
TAG=$(git rev-parse --short HEAD)
IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/${SERVICE}:${TAG}"

docker build --platform linux/amd64 -t "${IMAGE}" .
docker push "${IMAGE}"

gcloud run deploy "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image "${IMAGE}" \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars "WEATHER_API_KEY=...,WEATHER_API_URL=https://api.weatherapi.com/v1/current.json,VIACEP_API_URL=https://viacep.com.br/ws"
```

## Cloud Run runtime notes

- Cloud Run injects `PORT` at runtime; `internal/config.LoadConfig` honors it over
  `SERVER_PORT`. The container therefore needs no port-specific configuration.
- Images must be **linux/amd64** unless you opt into ARM — the `docker build` step in
  the script forces the platform explicitly so Apple Silicon hosts produce a compatible
  image.
- The service is deployed with `--allow-unauthenticated` so the grader can hit it
  directly.

## Hardening (recommended for production)

`WEATHER_API_KEY` is passed as a plain Cloud Run env var by default to keep the
challenge deploy reproducible in a single command. For production, store the key in
**Secret Manager** and inject it as a secret-backed env var:

```bash
echo -n "<api-key>" | gcloud secrets create weather-api-key --data-file=-

gcloud run services update weather-by-zipcode \
  --region us-central1 \
  --remove-env-vars WEATHER_API_KEY \
  --update-secrets WEATHER_API_KEY=weather-api-key:latest
```

The Cloud Run service account also needs the `roles/secretmanager.secretAccessor` role
on the secret.

## Verify

After the script finishes, hit each of the challenge scenarios against the printed URL:

```bash
URL=$(gcloud run services describe weather-by-zipcode \
  --region us-central1 --format='value(status.url)')

curl -i "$URL/health"                  # 200 ok
curl -i "$URL/weather/01001000"        # 200 + temp_C/temp_F/temp_K
curl -i "$URL/weather/123"             # 422 invalid zipcode
curl -i "$URL/weather/00000000"        # 404 can not find zipcode
```

Update the **Deployment** section of `README.md` with the resulting URL.

## Teardown

```bash
gcloud run services delete weather-by-zipcode --region us-central1
gcloud artifacts repositories delete weather-by-zipcode --location us-central1
```