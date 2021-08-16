FROM golang:1.16-buster
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install git
WORKDIR /work

all:
    BUILD +build
    BUILD +lint

build:
    ARG service
    RUN go get -u github.com/ahmetb/govvv
    RUN go get -u github.com/go-bindata/go-bindata/...
    COPY . ./
    RUN go generate ./pkg/$service/run.go
    RUN govvv build -pkg github.com/smartnuance/saas-kit/pkg/$service ./cmd/$service
    SAVE ARTIFACT $service AS LOCAL bin/$service

lint:
    RUN go get golang.org/x/lint/golint
    COPY cmd pkg ./
    RUN golint -set_exit_status ./...
