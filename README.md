# go SAAS kit

A reusable set of micro-services for a multi-tenant SAAS backend.

## Services & features

[Authentication & permission service](./pkg/auth):
- authenticates users and issues JWTs with appropriate roles/permissions
- supports revocation of JWTs
- multiple profiles per user, profiles belong to one SAAS instance

## A note on the frontend

This example repository is integrated with a [compatible frontend built with Flutter (for web)](https://github.com/smartnuance/flutter-admin-kit). It can be used with any frontend technology since the exposed services offer complete REST APIs.

## Getting started

### Environment

Choose your environment `dev` (default), `test` or `prod`:

> export SAAS_KIT_ENV="test"

Create an environment file `.env.dev` or `.env.test` for the environments used, where you might want to override some env-specific variables.

###  database

Create a local postgres instance or use the provided `docker-compose.yml` file to set it up:

> docker network create saas-kit-net

> docker volume create --name=saas-kit-pg

> docker-compose --env-file .env.dev up

### Go-bindata

For developers's convenience all commands necessary for development like code generation is done via `go:generate` commands at the top of each service's `run.go` or for very lazy developers they are all delegated from one file and invoked by a single generate command:

> go generate ./pkg/auth

### Run service(s)

Run all sevices by separate go routines:

> go run ./cmd/dev

Run single service:

> go run ./cmd/auth


### Build service(s)

To include build information we use the [`govvv`](github.com/ahmetb/govvv) utility:

> go get github.com/ahmetb/govvv

Then

> govvv build ./cmd/auth -pkg github.com/smartnuance/saas-kit/pkg/auth

## Packages used

API framework:
- [gin](https://github.com/gin-gonic/gin)

Database interaction:
- [gorm](https://gorm.io) for postgres, **not using auto-migrations**
- [golang-migrate](https://github.com/golang-migrate/migrate) for **clean up/down migrations**

Asset handling:
- [go-bindata](https://github.com/go-bindata/go-bindata)