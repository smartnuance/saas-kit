# go SaaS kit

A reusable set of micro-services for a multi-tenant SaaS backend.

<a href="https://smartnuance.com/blog/saas-kit/"><img src="https://gitlab.com/smartnuance/blog/-/raw/c0bcf702cabdcc432449d816da401a7f3e1fd6b6/src/data/blog/saas-kit/saas-kit.svg?inline=false" alt="go SaaS kit" width="80"></a>

[demo](https://admin-kit.smartnuance.com) | [introductory blog post](https://smartnuance.com/blog/saas-kit/) | [API playground](https://mybinder.org/v2/gh/smartnuance/saas-kit/HEAD?labpath=test%2Fnotebooks%2Fapi.ipynb)


## Services & features

- Authentication & Authorization: [service](./pkg/auth) | [design](https://smartnuance.com/blog/saas-kit/#token-authentication) | [API playground](https://mybinder.org/v2/gh/smartnuance/saas-kit/HEAD?labpath=test%2Fnotebooks%2Fapi.ipynb)


## A note on the frontend

This project is integrated with a [compatible frontend built with Flutter (for web)](https://github.com/smartnuance/flutter-admin-kit). It can be used with any frontend technology since the exposed services offer complete REST APIs.


## Getting started

### Environment

Choose your environment `dev` (default), `test` or `prod`:

> export SAAS_KIT_ENV="test"

Create an environment file `.env.dev` or `.env.test` for the environments used, where you might want to override some env-specific variables.


### Database

Create a local postgres instance or use the provided `docker-compose.yml` file to set it up:

> docker network create saas-kit-net

> docker volume create --name=saas-kit-pg

> docker-compose --env-file .env.dev up

Initialize separate database schemas, one for each service (in production you have service-local databases with only one schema):

> go run ./cmd/dev init

Migrate services:

> go run ./cmd/auth migrate

> go run ./cmd/event migrate


### Run service(s)

Run single service

> go run ./cmd/auth


Run all sevices by separate go routines
(migrations are not applied!):

> go run ./cmd/dev


### Create some necessary data

Load fixtures from yaml file with a default instance:

> go run ./cmd/auth fixture ./pkg/auth/fixtures/instances.yml

Create super user:

> go run ./cmd/auth adduser -name=Simon -email=simon@smartnuance.com -password=f00bartest -instance=smartnuance.com


### Interact with auth service

Use exposed endpoint to signup user (with no permissions to start with):

> http PUT :8801/signup instance="smartnuance.com" name=Bob email=bob@smartnuance.com password=alice

Test login and save refresh/access tokens:

> RES=$(http POST :8801/login email=simon@smartnuance.com password=f00bartest instance=smartnuance.com -v -b)

> RT=$(echo $RES | jq -r '.refreshToken') 

> AT=$(echo $RES | jq -r '.accessToken')


Refresh token:

> http -v POST :8801/refresh refreshToken=$RT

Revoke token:

> http -v DELETE :8801/revoke/ Authorization:"Bearer $AT"

Revoke all tokens:

> http -v DELETE :8801/revoke/all Authorization:"Bearer $AT"

For a while, I can still use the access token, for example to rerun the idem-potent revoke:

> http -v DELETE :8801/revoke/ Authorization:"Bearer $AT"

But if we try to use the revoked refresh token in a refresh call, this will fail:

> http -v POST :8801/refresh refreshToken=$RT


### Interact with event service

Since no implicit switch from the super admin is allowed, we provide the role header to temporarily switch to the _event organizer_ role:

> http -v PUT :8802/workshop Authorization:"Bearer $AT" role:"event organizer" instance:"c5263570ono4ui8qfhgg" title=Bachata locationName=Ponto


## Packages used

API framework:
- [gin](https://github.com/gin-gonic/gin)

Database interaction:
- [sqlboiler](https://github.com/volatiletech/sqlboiler#getting-started) for generated db interaction
- [golang-migrate](https://github.com/golang-migrate/migrate) for **clean up/down migrations**
- [Globally Unique ID Generator]("github.com/rs/xid") that uses Mongo Object ID algorithm

Token handling:
- [jwt-go](https://github.com/golang-jwt/jwt) (v4!)

Logging:
- [zerolog](https://github.com/rs/zerolog)

Asset handling:
- [go-bindata](https://github.com/go-bindata/go-bindata)

Environment & Building:
- [godotenv](https://github.com/joho/godotenv)
- [govvv](https://github.com/ahmetb/govvv)
- [EARTHLY](https://docs.earthly.dev)

Testing:
- [gomock](https://github.com/golang/mock)
- [go-testdep](https://github.com/maxatome/go-testdeep)


## Contribute

### Development tools

All tools necessary for development like installing code generators is done via `go:generate` commands at the top of `cmd/dev/main.go`:

> go generate ./cmd/dev

(a subset of those tools are installed in CI build step of `Earthfile`)

For developers's convenience all generation commands are collected via `go:generate` at the top of each service's `run.go`, e.g. for the auth service

> go generate ./pkg/auth


### Migrate database

To interact with database, we use a schema first approach with [sqlboiler](https://github.com/volatiletech/sqlboiler#getting-started). It generates type-safe code to interact with the DB.

To start from an empty database (and test down migrations):

> go run ./cmd/auth migrate -down

Or if you messed the database up and want to start from scratch, you can reinit the database (**data will be lost!**):

> go run ./cmd/dev deinit

> go run ./cmd/dev init

Migrate services:

> go run ./cmd/auth migrate

> go run ./cmd/event migrate

When database is on newest version, we have to generated git-versioned DB models by

> go generate ./pkg/auth/db.go


### Build service(s)

To include build information we use the [`govvv`](github.com/ahmetb/govvv) utility:

> go install github.com/ahmetb/govvv@latest

Then

> govvv build -pkg github.com/smartnuance/saas-kit/pkg/auth -o ./bin/ ./cmd/auth

To create reproducable builds, you can use [EARTHLY](https://docs.earthly.dev):

> earthly --build-arg service=auth +build

With either way the resulting runnable is executed by:

> ./bin/auth


### Deploy services

Again with [EARTHLY](https://docs.earthly.dev), we can switch the target to `deploy`:

> earthly --build-arg service=dev +publish

This will not yet push the created single-binary docker image.

You can try to run the docker image locally, using the host network (so the ports used by the services needs to be free):

> docker run --rm --network host -v $PWD/test/data:/app/test/data ghcr.io/smartnuance/saas-kit:latest


### Configure linter

We are using the configurable [staticcheck](https://staticcheck.io/docs/) linter.

> staticcheck ./...

Check that your [VScode](https://code.visualstudio.com/) workspace settings contain

    "go.lintTool": "staticcheck"


### Configure testing

If your IDE (like VScode) shows broken test output because it does not support colored output, add this to your workspace settings:

    "go.testEnvVars": {
        "TESTDEEP_COLOR": "off"
    },
