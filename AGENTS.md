# Repository Guidelines

## Project Structure & Module Organization
Use a split layout:
- `web/`: Vue 3 + TypeScript + Ant Design Vue frontend.
- `server/`: Go backend with Gin and Gorm.
- `docs/`: API notes, ER diagrams, UI conventions, and deployment docs.

Frontend code should live under `web/src`, with `views/`, `components/`, `api/`, `stores/`, and `router/`. Backend code should follow `server/cmd/`, `server/internal/handler`, `server/internal/service`, `server/internal/repository`, `server/internal/model`, and `server/internal/dto`.

## Build, Test, and Development Commands
Prefer explicit module-level commands:
- `cd web && pnpm dev`: start the frontend dev server.
- `cd web && pnpm build`: produce the production bundle.
- `cd web && pnpm lint && pnpm type-check`: enforce TS and Vue quality gates.
- `cd server && go run ./cmd/api`: run the Gin service locally.
- `cd server && go test ./...`: run backend tests.
- `cd server && go vet ./...`: catch common Go issues.

Keep scripts or a `Makefile` as thin wrappers around these commands.

## Coding Style & Naming Conventions
Frontend uses `script setup` with TypeScript. Name Vue components in PascalCase, composables as `useXxx`, Pinia stores as `useXxxStore`, and API modules by domain such as `user.ts` or `auth.ts`. Keep page files focused on composition; move request and transform logic into `api/` or composables. Prefer Ant Design Vue components before building custom widgets.

Backend code must be `gofmt`-formatted. Keep package names lowercase and short. Handlers only bind, validate, and return responses. Business rules belong in `service`, and Gorm access belongs in `repository`. Do not return database models directly; define DTO or VO structs explicitly.

## Testing Guidelines
Frontend tests should cover composables, state transitions, and critical page behavior once the test stack is added. Backend tests should be table-driven where practical and named `*_test.go`. Add a regression test for bug fixes when feasible. Before opening a PR, run frontend lint and type-check plus `go test ./...`.

## Commit & Pull Request Guidelines
Use short imperative commit subjects, for example `Add login page scaffold` or `Refactor user repository`. Keep backend and frontend changes separate when they are logically independent. PRs should include the summary, affected paths, verification commands, linked issue if available, and screenshots for visible UI changes.

## Security & Configuration Tips
Do not commit `.env`, local secrets, IDE metadata, build output, or coverage artifacts. Keep configuration in example files such as `.env.example`, and document required variables in `docs/` rather than storing credentials in source.

## Agent Response Language
Use Simplified Chinese for all user-visible content by default, including progress updates, plans, confirmation prompts, and final replies. Do not expose full chain-of-thought; provide short Chinese reasoning summaries when explanation is useful. Switch languages only if the user explicitly requests it.
