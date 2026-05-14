# syntax=docker/dockerfile:1

FROM golang:1.24-bookworm AS build
WORKDIR /src
ENV GOPROXY=direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/kindling .

FROM python:3.12-slim-bookworm AS runtime

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates imagemagick \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /out/kindling /app/kindling
COPY Fondamento-Regular.ttf SourceCodePro-Regular.ttf WinterSong-owRGB.ttf impact.ttf /app/
COPY scripts/kindling-terminus-loop.py /app/kindling-terminus-loop.py
RUN mkdir -p /app/trmnl /uploads /state \
    && chown -R 1000:1000 /app/trmnl /uploads /state \
    && chmod +x /app/kindling /app/kindling-terminus-loop.py

ENV APP_DIR=/app \
    KINDLING_BIN=/app/kindling \
    SOURCE_IMAGE=/app/trmnl/weather.png \
    STAGING_DIR_LOCAL=/uploads \
    STAGING_FILENAME=kindling_weather.png \
    STAGING_PATH_TERMINUS=/app/public/uploads/kindling_weather.png \
    STATE_FILE=/state/kindling-weather-screen-id \
    INTERVAL_SECONDS=300

CMD ["/app/kindling-terminus-loop.py"]
