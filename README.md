# Kindling

Kindling renders an 800×480 e-ink weather image for TRMNL / Kindle-style displays.

The current weather screen shows Brisbane and Saint-Émilion conditions from
OpenWeatherMap, plus a timestamp so it is obvious when the display stopped
refreshing.

## Local run

```sh
export OPENWEATHERMAP_API_KEY=...
go run .
```

Outputs land under `trmnl/`, including `trmnl/weather.png`.

## Containerized Terminus mode

The preferred homelab deployment is a Docker Compose stack on `tatajuba` at:

```text
/home/drone/staging/tatajuba/kindling
```

The container:

1. runs the Go renderer;
2. converts `trmnl/weather.png` to a 2-bit 800×480 PNG with ImageMagick;
3. writes the PNG into the shared `trmnl-byos_web-uploads` Docker volume; and
4. POSTs/PATCHes Terminus with `uri + preprocessed: true`.

This mirrors the `photo-gallery` TRMNL integration and avoids Terminus's HTML
sanitizer path, which strips `data:` image URIs.

Deployment shape:

```sh
cp .env.example .env
# fill OPENWEATHERMAP_API_KEY, TERMINUS_LOGIN, TERMINUS_PASSWORD

docker compose up -d --build
```

The default compose settings attach the screen to playlist `3`, matching the
current OG TRMNL apartment playlist.

## Legacy tome mode

The old `systemd/` unit and `systemd/scripts/fire.sh` were used on the `tome`
blog server to publish images under `https://sploo.sh/trmnl/`. They are retained
for reference, but the intended operational model is now the containerized
Terminus publisher above.
