#!/usr/bin/env bash
# Build the container image, push to Google Artifact Registry, and
# deploy the weather service to Cloud Run.
#
# Required env vars:
#   PROJECT_ID        Google Cloud project ID.
#   WEATHER_API_KEY   WeatherAPI.com API key.
#
# Optional env vars (defaults shown):
#   REGION            us-central1
#   SERVICE           weather-by-zipcode
#   REPO              weather-by-zipcode
#   IMAGE_TAG         git short SHA, or "latest" outside a git checkout
#   WEATHER_API_URL   https://api.weatherapi.com/v1/current.json
#   VIACEP_API_URL    https://viacep.com.br/ws
#
# One-time prerequisites are documented in docs/DEPLOY.md.

set -euo pipefail

: "${PROJECT_ID:?PROJECT_ID is required}"
: "${WEATHER_API_KEY:?WEATHER_API_KEY is required}"

REGION="${REGION:-us-central1}"
SERVICE="${SERVICE:-weather-by-zipcode}"
REPO="${REPO:-weather-by-zipcode}"
WEATHER_API_URL="${WEATHER_API_URL:-https://api.weatherapi.com/v1/current.json}"
VIACEP_API_URL="${VIACEP_API_URL:-https://viacep.com.br/ws}"

if [[ -z "${IMAGE_TAG:-}" ]]; then
  if git rev-parse --short HEAD >/dev/null 2>&1; then
    IMAGE_TAG="$(git rev-parse --short HEAD)"
  else
    IMAGE_TAG="latest"
  fi
fi

IMAGE="${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPO}/${SERVICE}:${IMAGE_TAG}"

echo ">> Building image ${IMAGE}"
docker build --platform linux/amd64 -t "${IMAGE}" .

echo ">> Pushing image to Artifact Registry"
docker push "${IMAGE}"

echo ">> Deploying ${SERVICE} to Cloud Run (${REGION})"
gcloud run deploy "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image "${IMAGE}" \
  --platform managed \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars "WEATHER_API_KEY=${WEATHER_API_KEY},WEATHER_API_URL=${WEATHER_API_URL},VIACEP_API_URL=${VIACEP_API_URL}"

echo ">> Service URL:"
gcloud run services describe "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --format='value(status.url)'