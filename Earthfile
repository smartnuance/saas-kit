VERSION 0.6
FROM golang:1.18-bullseye
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install git
WORKDIR /work

all:
    BUILD +build
    BUILD +lint

setup:
    COPY . ./

    # ignore git-crypt filters because we do not necessarily have access to secrets encrypted in repo
    RUN git config --unset-all filter.git-crypt.smudge
    RUN git config --unset-all filter.git-crypt.clean
    RUN git config --unset-all filter.git-crypt.required
    RUN git config --unset-all diff.git-crypt.textconv

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
    COPY .env* .
    COPY pkg/event/modelinfo ./modelinfo

    EXPOSE 8800
    ENTRYPOINT ["/app/dev"]

    SAVE IMAGE --push ghcr.io/smartnuance/saas-kit:latest

lint:
    FROM +setup
    RUN go install honnef.co/go/tools/cmd/staticcheck@2022.1

    RUN staticcheck ./...
