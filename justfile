
APP := "pomogo"


generate:
    @go generate ./cmd/pomogo...

build:
    @go build -o ./{{APP}} ./cmd/pomogo...

test:
    @go test -coverprofile=c.out ./... -json

coverage-html:
    [ -f c.out ] || >&2 echo "Must run test first"
    [ -f c.out ] && go tool cover -html=c.out

coverage:
    [ -f c.out ] || >&2 echo "Must run test first"
    [ -f c.out ] && go tool cover -func=c.out

lint:
    golangci-lint run ./...

fmt:
    go fmt ./...

clean:
    rm -f c.out
    go clean


pre-commit: generate fmt
