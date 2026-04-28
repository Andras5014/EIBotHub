# EIBotHub Deploy Package

Release artifacts contain one executable and this `deploy/` directory.

## Files

- `config.json`: runtime configuration loaded by the executable when present.
- `start-linux.sh`: starts the Linux binary.
- `start-windows.ps1`: starts the Windows binary.

## Config

The executable looks for `deploy/config.json` first. Relative `db_path` and
`storage_dir` values in `deploy/config.json` are resolved from the release root.
Environment variables already set on the host override values in `config.json`.

- `port` / `APP_PORT`: HTTP port, default `8080`.
- `db_path` / `DB_PATH`: SQLite database path.
- `storage_dir` / `STORAGE_DIR`: uploaded file storage path.
- `app_secret` / `APP_SECRET`: JWT and download-token secret. Change this before production use.
- `seed_demo` / `SEED_DEMO`: set to `false` to disable demo seed data.
- `gin_mode` / `GIN_MODE`: use `release` for packaged deployments.

Set `CONFIG_FILE` to load a custom JSON file from another path.
