# Instructions for contributors

* Run `go test ./...` and ensure all tests pass before committing.
* Verify cross compilation for macOS arm64:
  `GOOS=darwin GOARCH=arm64 go build ./cmd/dpxe` and `GOOS=darwin GOARCH=arm64 go test ./...`.
* Use `logrus` for logging and keep code formatted with `gofmt -w`.
