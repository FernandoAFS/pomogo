
APP := "pomogo"

build:
    @go build -o ./{{APP}} ./cmd...

test:
    @go test -coverprofile=c.out ./... -json

coverage-html:
    [ -f c.out ] || >&2 echo "Must run test first"
    [ -f c.out ] && go tool cover -html=c.out

coverage:
    [ -f c.out ] || >&2 echo "Must run test first"
    [ -f c.out ] && go tool cover -func=c.out

clean:
    rm -f c.out
    go clean

