# AGENT Instructions

- Use `go test ./...` before committing any changes to verify code builds and tests pass.
- Keep code in idiomatic Go style using modules under `internal/` and command entry in `cmd/dpxe`.
- Logging should use logrus with full timestamps. HTTP logs must mimic nginx combined style.
