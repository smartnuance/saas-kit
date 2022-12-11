VERSION 0.6
FROM golang:1.18-bullseye
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install git git-crypt
WORKDIR /work

all:
    BUILD +build
    BUILD +lint

setup:
    COPY . ./

build:
    FROM +setup
    ARG service
    RUN go install github.com/ahmetb/govvv@latest
    
    RUN govvv build ./cmd/$service
    SAVE ARTIFACT $service AS LOCAL bin/$service

publish:
    FROM debian:bullseye

    WORKDIR /app
    COPY +build/$service .
    COPY .env .
    COPY public ./public
    COPY ./pkg/auth/fixtures/instances.yml .

    EXPOSE 8800
    ENTRYPOINT ["/app/dev"]

    SAVE IMAGE --push ghcr.io/smartnuance/saas-kit:latest

lint:
    FROM +setup
    RUN go install honnef.co/go/tools/cmd/staticcheck@2022.1

    RUN staticcheck ./...
