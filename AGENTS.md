# Repository Guidelines

## Project Structure & Module Organization

- `kv/`: TinyKV server, storage engine, Raftstore, transaction logic, and related tests (including `kv/test_raftstore`).
- `raft/`: Raft consensus implementation and its unit tests.
- `scheduler/`: TinyScheduler server, client, configs, and end-to-end tests under `scheduler/tests`.
- `proto/`: Protocol Buffers definitions and generated Go code; update via `make proto` when `.proto` files change.
- `log/`, `doc/`, `scripts/`: Logging utilities, course documentation, and CI/performance scripts.

## Build, Test, and Development Commands

- `make` or `make kv scheduler`: Build `tinykv-server` and `tinyscheduler-server` into `bin/`.
- `make test`: Run Go unit/integration tests with coverage for most packages.
- `make project1|project2|project3|project4`: Run project-specific grading suites for each course stage.
- `make proto`: Regenerate protobuf Go code after modifying definitions.
- `make format` / `make ci`: Format all Go files and run `gofmt` + `go vet` checks.

## Coding Style & Naming Conventions

- Go 1.13+; always run `gofmt` (or `make format`) before committing.
- Use idiomatic Go naming: `CamelCase` for exported symbols, `lowerCamel` for unexported; keep names descriptive and consistent with existing code.
- Keep packages cohesive (e.g., `kv/storage`, `kv/transaction`); avoid large, mixed-purpose files.
- Use the logging utilities in `log/` and `github.com/pingcap/log` instead of `fmt.Println`.

## Testing Guidelines

- Add or update tests next to the code you change (e.g., `raft/*_test.go`, `kv/test_raftstore`, `scheduler/tests`).
- Ensure `make test` and any relevant `make projectN` targets pass before opening a PR.
- Name tests clearly to reflect behavior and project stage, e.g., `TestCommitMissingPrewrite4B`.

## Commit & Pull Request Guidelines

- Use concise, imperative commit messages, optionally with scope, e.g., `raft: fix leader election` or `transaction: handle missing lock (#123)`.
- In PR descriptions, include: problem statement, overview of changes, important design choices, and how you tested (commands + key outputs).
- Link related GitHub issues or course tasks; attach logs or screenshots only when they clarify failures or performance changes.

## Agent-Specific Instructions

- Keep edits minimal and localized to relevant packages; do not refactor broadly without clear motivation.
- Prefer validating changes with `make test` and the smallest relevant `make projectN` target.
- Avoid introducing new dependencies unless necessary; reuse existing patterns and libraries in this repository.

