FROM golang:1.17-buster
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install git
WORKDIR /work

all:
    BUILD +build
    BUILD +lint

build:
    ARG service
    RUN go install github.com/ahmetb/govvv@latest
    
    COPY . ./
    RUN govvv build ./cmd/$service
    SAVE ARTIFACT $service AS LOCAL bin/$service

publish:
    FROM debian:buster

    WORKDIR /app
    COPY +build/$service .
    COPY .env .
    COPY .env* .
    COPY pkg/event/modelinfo ./modelinfo
    COPY prod ./prod

    EXPOSE 8800
    ENTRYPOINT ["/app/dev"]

    SAVE IMAGE --push ghcr.io/smartnuance/saas-kit:latest

lint:
    RUN go install honnef.co/go/tools/cmd/staticcheck@2021.1
    COPY . ./
    RUN staticcheck ./...
