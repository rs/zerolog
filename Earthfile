ARG GO_VERSION=1.16
FROM golang:$GO_VERSION-buster

deps:
    RUN apt update && apt install -y systemd
    RUN go get golang.org/x/tools/cmd/goimports
    RUN go get golang.org/x/lint/golint
    RUN go get github.com/gordonklaus/ineffassign

code:
    FROM +deps
    WORKDIR /zerolog
    COPY --dir cmd diode hlog internal journald log pkgerrors go.mod go.sum .
    RUN go mod download
    COPY *.go .

lint:
    FROM +code
    RUN output="$(ineffassign ./... 2>&1 | grep -v '/earthly/ast/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d $(find . -type f -name '*.go' | grep -v \.pb\.go) 2>&1)"  ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN golint -set_exit_status ./...
    RUN output="$(go vet ./... 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi

test:
    FROM +code
    RUN /lib/systemd/systemd-journald & \
        go test -race -cpu=1,2,4 -bench . -benchmem ./... && \
        go test -tags binary_log -race -cpu=1,2,4 -bench . -benchmem ./...

    # TODO expand on this test to validate the correct logs are received; currently
    # there's a TODO for this under
    # https://github.com/rs/zerolog/blob/master/journald/journald_test.go#L22-L50
    RUN journalctl -o verbose | grep 'JSON={"level":"info","message":"Tick!"}'

all:
    # Currently the linter is failing, but you can enable it here:
    #BUILD +lint

    BUILD \
        --build-arg GO_VERSION=1.16 \
        --build-arg GO_VERSION=1.15 \
        --build-arg GO_VERSION=1.14 \
        --build-arg GO_VERSION=1.13 \
        --build-arg GO_VERSION=1.12 \
        +test
