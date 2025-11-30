#!/usr/bin/env bash

set -euo pipefail

# Defaults
CHECK_TIMEOUT="${CHECK_TIMEOUT:-3}"
CHECK_HTTP_CODE="${CHECK_HTTP_CODE:-200}"
RETRY_COUNT="${RETRY_COUNT:-3}"
RETRY_DELAY="${RETRY_DELAY:-1}"

# Validation
: "${SERVICE_NAME:?SERVICE_NAME is required: the systemd service to restart}"
: "${CHECK_URL:?CHECK_URL is required: the healthcheck URL}"
: "${CHECK_HTTP_CODE:?CHECK_HTTP_CODE is required: expected HTTP status code}"
: "${CHECK_TIMEOUT:?CHECK_TIMEOUT is required: curl timeout in seconds}"
: "${RETRY_COUNT:?RETRY_COUNT is required: number of retry attempts}"
: "${RETRY_DELAY:?RETRY_DELAY is required: delay between retry attempts}"

log() {
    echo $(date +"%Y-%m-%d %H:%M:%S") $1
}

check_app() {
    curl -s -o /dev/null -w "%{http_code}" --max-time "$CHECK_TIMEOUT" "$CHECK_URL"
}

log "INFO: SERVICE_NAME=$SERVICE_NAME \
    CHECK_URL=$CHECK_URL \
    CHECK_HTTP_CODE=$CHECK_HTTP_CODE \
    CHECK_TIMEOUT=$CHECK_TIMEOUT \
    RETRY_COUNT=$RETRY_COUNT \
    RETRY_DELAY=$RETRY_DELAY"

failures=0

for ((i=0; i<RETRY_COUNT; i++)); do
    code=$(check_app || echo "failed")

    if [[ "$code" == "$CHECK_HTTP_CODE" ]]; then
        log "INFO: healthcheck returned $code"
        exit 0
    else
        log "WARN: healthcheck failed (code=$code), attempt $i/$RETRY_COUNT"
        failures=$((failures+1))
        sleep "$RETRY_DELAY"
    fi
done

log "ERROR: healthcheck failed $failures times, restarting service $SERVICE_NAME"

if systemctl restart "$SERVICE_NAME"; then
    log "INFO: service $SERVICE_NAME restarted successfully"
else
    log "ERROR: failed to restart service $SERVICE_NAME"
fi
