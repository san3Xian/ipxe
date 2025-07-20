# Guidelines for contributors

- Run `gofmt -w` on all Go files before committing.
- Ensure tests pass with:
  - `GOOS=darwin GOARCH=arm64 go test ./...`
  - `GOOS=darwin GOARCH=arm64 go build ./cmd/dpxe`
- Use logrus for logging and keep HTTP request log format similar to nginx.
- Keep modules within `internal/` directory for project code.
- Add descriptive comments where helpful.

