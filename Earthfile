FROM golang:1.16-buster
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install git
WORKDIR /work

all:
    BUILD +build
    BUILD +lint

build:
    ARG service
    RUN go install github.com/ahmetb/govvv@latest
    
    COPY . ./
    RUN go generate ./pkg/$service/run.go
    RUN govvv build -pkg github.com/smartnuance/saas-kit/pkg/$service ./cmd/$service
    SAVE ARTIFACT $service AS LOCAL bin/$service

lint:
    RUN go install honnef.co/go/tools/cmd/staticcheck@2021.1
    COPY . ./
    RUN staticcheck ./...
