package sourcegraph

// Keep these versions in sync with go.mod
//go:generate env GOBIN=$PWD/.bin GO111MODULE=on go install golang.org/x/tools/cmd/goimports@v0.17.0
//go:generate go run github.com/derision-test/go-mockgen/cmd/go-mockgen@v1.3.8-0.20240105000756-fb9effb23d90
