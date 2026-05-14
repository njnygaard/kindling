#!/usr/bin/env python3
"""Render Kindling weather and publish it to Terminus on a loop.

The container runs the Go Kindling binary to generate /app/trmnl/weather.png,
ImageMagick converts it into a 2-bit 800x480 PNG in the shared Terminus uploads
volume, then the Terminus screens API is POSTed/PATCHed with uri+preprocessed.
"""

from __future__ import annotations

import json
import os
import subprocess
import sys
import time
import urllib.error
import urllib.request
from datetime import datetime
from pathlib import Path

APP_DIR = Path(os.environ.get("APP_DIR", "/app"))
KINDLING_BIN = Path(os.environ.get("KINDLING_BIN", "/app/kindling"))
SOURCE_IMAGE = Path(os.environ.get("SOURCE_IMAGE", "/app/trmnl/weather.png"))
STAGING_DIR_LOCAL = Path(os.environ.get("STAGING_DIR_LOCAL", "/uploads"))
STAGING_FILENAME = os.environ.get("STAGING_FILENAME", "kindling_weather.png")
STAGING_PATH_TERMINUS = os.environ.get(
    "STAGING_PATH_TERMINUS", f"/app/public/uploads/{STAGING_FILENAME}"
)
CONVERT_BIN = os.environ.get("CONVERT_BIN", "convert")
TERMINUS_URL = os.environ.get("TERMINUS_URL", "http://host.docker.internal:2300")
TERMINUS_LOGIN = os.environ.get("TERMINUS_LOGIN", "")
TERMINUS_PASSWORD = os.environ.get("TERMINUS_PASSWORD", "")
SCREEN_NAME = os.environ.get("SCREEN_NAME", "kindling_weather")
SCREEN_LABEL = os.environ.get("SCREEN_LABEL", "Kindling Weather")
MODEL_ID = int(os.environ.get("MODEL_ID", "1"))
PLAYLIST_ID = int(os.environ.get("PLAYLIST_ID", "3"))
INTERVAL_SECONDS = int(os.environ.get("INTERVAL_SECONDS", "300"))
STATE_FILE = Path(os.environ.get("STATE_FILE", "/state/kindling-weather-screen-id"))


def log(msg: str) -> None:
    sys.stderr.write(f"[{datetime.now().isoformat(timespec='seconds')}] kindling: {msg}\n")
    sys.stderr.flush()


def require_env(name: str) -> str:
    value = os.environ.get(name, "")
    if not value:
        raise RuntimeError(f"{name} is required")
    return value


def run_kindling() -> None:
    # The Go binary loads fonts and writes trmnl/*.png relative to cwd.
    subprocess.run([str(KINDLING_BIN)], cwd=APP_DIR, check=True)
    if not SOURCE_IMAGE.exists():
        raise RuntimeError(f"Kindling did not produce {SOURCE_IMAGE}")


def convert_to_staging() -> int:
    STAGING_DIR_LOCAL.mkdir(parents=True, exist_ok=True)
    staging_local = STAGING_DIR_LOCAL / Path(STAGING_FILENAME).name
    subprocess.run(
        [
            CONVERT_BIN,
            str(SOURCE_IMAGE),
            "-colorspace",
            "Gray",
            "-dither",
            "FloydSteinberg",
            "-posterize",
            "4",
            "-alpha",
            "off",
            "-depth",
            "2",
            "-strip",
            f"PNG:{staging_local}",
        ],
        check=True,
    )
    return staging_local.stat().st_size


def http(method: str, path: str, *, body: dict | None = None, token: str | None = None) -> dict:
    url = f"{TERMINUS_URL}{path}"
    data = json.dumps(body).encode() if body is not None else None
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    if token:
        headers["Authorization"] = token
    req = urllib.request.Request(url, data=data, method=method, headers=headers)
    with urllib.request.urlopen(req, timeout=30) as resp:
        text = resp.read().decode()
        return json.loads(text) if text else {}


def login() -> str:
    resp = http(
        "POST",
        "/login",
        body={"login": require_env("TERMINUS_LOGIN"), "password": require_env("TERMINUS_PASSWORD")},
    )
    token = resp.get("access_token")
    if not token:
        raise RuntimeError(f"login response missing access_token: {resp}")
    return token


def find_existing_screen_id(token: str) -> int | None:
    try:
        existing = http("GET", "/api/screens", token=token)
    except urllib.error.HTTPError as e:
        log(f"screen list failed: HTTP {e.code} {e.reason}")
        return None
    for screen in existing.get("data") or []:
        if screen.get("name") == SCREEN_NAME and screen.get("model_id") == MODEL_ID:
            return int(screen["id"])
    return None


def attach_to_playlist(token: str, screen_id: int) -> None:
    try:
        playlist = http("GET", f"/api/playlists/{PLAYLIST_ID}", token=token)
    except urllib.error.HTTPError as e:
        log(f"playlist get failed: HTTP {e.code} {e.reason} — skipping attach")
        return

    pl = playlist.get("data") or playlist
    items = pl.get("items") or []
    if any(item.get("screen_id") == screen_id for item in items):
        log(f"screen {screen_id} already on playlist {PLAYLIST_ID}")
        return

    rebuilt = [
        {"screen_id": item["screen_id"], "position": idx + 1}
        for idx, item in enumerate(items)
    ] + [{"screen_id": screen_id, "position": len(items) + 1}]
    body = {
        "playlist": {
            "name": pl.get("name", "default"),
            "label": pl.get("label", "Default"),
            "mode": pl.get("mode", "automatic"),
            "items": rebuilt,
        }
    }
    try:
        http("PATCH", f"/api/playlists/{PLAYLIST_ID}", body=body, token=token)
        log(f"attached screen {screen_id} to playlist {PLAYLIST_ID} (now {len(rebuilt)} items)")
    except urllib.error.HTTPError as e:
        detail = e.read().decode(errors="replace")[:300]
        log(f"playlist attach failed: HTTP {e.code} {e.reason} | {detail}")


def push_screen(token: str, screen_id: int | None) -> int | None:
    payload = {
        "screen": {
            "name": SCREEN_NAME,
            "label": SCREEN_LABEL,
            "model_id": MODEL_ID,
            "uri": STAGING_PATH_TERMINUS,
            "preprocessed": True,
        }
    }
    if screen_id is None:
        try:
            create = http("POST", "/api/screens", body=payload, token=token)
        except urllib.error.HTTPError as e:
            detail = e.read().decode(errors="replace")[:300]
            log(f"POST failed: HTTP {e.code} {e.reason} | {detail}")
            return None
        screen_id = (create.get("data") or create).get("id")
        log(f"POST /api/screens created id={screen_id}")
        return int(screen_id) if screen_id else None

    try:
        http("PATCH", f"/api/screens/{screen_id}", body=payload, token=token)
        log(f"PATCH /api/screens/{screen_id} ok")
    except urllib.error.HTTPError as e:
        detail = e.read().decode(errors="replace")[:300]
        log(f"PATCH failed: HTTP {e.code} {e.reason} | {detail}")
    return screen_id


def tick() -> None:
    run_kindling()
    bytes_written = convert_to_staging()
    log(f"wrote {bytes_written} bytes to {STAGING_PATH_TERMINUS}")

    token = login()

    screen_id: int | None = None
    if STATE_FILE.exists() and STATE_FILE.read_text().strip():
        try:
            screen_id = int(STATE_FILE.read_text().strip())
        except ValueError:
            screen_id = None
    if screen_id is None:
        screen_id = find_existing_screen_id(token)
        if screen_id:
            STATE_FILE.parent.mkdir(parents=True, exist_ok=True)
            STATE_FILE.write_text(str(screen_id))
            log(f"recovered existing screen id {screen_id} from API")

    new_id = push_screen(token, screen_id)
    if new_id:
        STATE_FILE.parent.mkdir(parents=True, exist_ok=True)
        STATE_FILE.write_text(str(new_id))
        attach_to_playlist(token, new_id)


def main() -> int:
    require_env("OPENWEATHERMAP_API_KEY")
    require_env("TERMINUS_LOGIN")
    require_env("TERMINUS_PASSWORD")
    STATE_FILE.parent.mkdir(parents=True, exist_ok=True)
    log(
        f"starting; source={SOURCE_IMAGE} staging={STAGING_PATH_TERMINUS} "
        f"terminus={TERMINUS_URL} interval={INTERVAL_SECONDS}s playlist={PLAYLIST_ID}"
    )
    while True:
        try:
            tick()
        except Exception as e:  # keep the service alive; the next tick may recover
            log(f"tick raised {type(e).__name__}: {e}")
        time.sleep(INTERVAL_SECONDS)


if __name__ == "__main__":
    sys.exit(main())
