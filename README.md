# Introduction

This repo is the Marble backend implementation:

- 1 single go app
- Postgres DB

## Getting Started

### Setup your environment

[Install Go](https://go.dev/doc/install) on your laptop (see the version in go.mod).

> NB: To handle different versions, you can look at [Managing Go installations](https://go.dev/doc/manage-install) or use a version manager tool like [asdf](https://github.com/kennyp/asdf-golang)

Create you own `.env` file based on `.env.example`. You can fill it with your own values but it should work locally with the default values (except for third-party services functionalities).

#### Setup the DB

Launch the postgres DB used by the backend:

```sh
docker compose up -d
```

> NB: it creates a `marble-backend_postgres-db` volume to store PG data.

#### Firebase emulator suite for local development

Install the [Firebase tools](https://firebase.google.com/docs/emulator-suite):

```sh
curl -sL firebase.tools | bash
```

Then copy the `./firebase-local-data.example` folder to `./firebase-local-data`. This folder will be used to store the local data of the Firebase emulator. It is ignored by git.

Then start it using (replace `GOOGLE_CLOUD_PROJECT` with the value from your `.env` file):

```sh
firebase --project GOOGLE_CLOUD_PROJECT emulators:start --import=./firebase-local-data --export-on-exit
```

> NB: The `--import` flag is used to import the local data into the emulator. The `--export-on-exit` flag is used to export the data when the emulator is stopped so you don't lose your changes.

### Launch the project

Export your `.env` file and run the root of the project:

```sh
go run . --migrations --server
```

> Using VSCode, you can also run the `Migrate and Launch (.env)` task in the "Run and debug" tab. This will load your env file, migrate the DB and start the server.

## Application flags

The application can be run with the following flags:

- `--migrations`: run the migrations
- `--server`: run the server
  - in development, it also runs the `SeedZorgOrganization` usecase script from `usecases/seed_usecase`
- `--scheduler`: run the scenario scheduling job
- `--scheduled-executer`: execute scheduled scenario job
- `--data-ingestion`: run data ingestion job

> NB: `.vscode/launch.json` contains the configuration to run the app with these flags.

## API

The rooting of the application is defined inside `api/routes.go`.

For further information on the API, you can also refer to the following resources:

- [our API docs](https://docs.checkmarble.com/reference/introduction-1) for public facing reference
- the Open API Specification defined in the frontend repository [here](https://github.com/checkmarble/marble-frontend/blob/main/packages/marble-api/scripts/openapi.yaml).

## DB Seed and migrations

The application uses [goose](https://github.com/pressly/goose) to manage migrations.

Migrations are located in the `repositories/migrations` folder.

Execute the program with flags `-migrations` to run migrations

To create a new migration, you can use the following command:

```sh
make generate_migration
```

## FAQ

### How to update firebase local data ?

- Run firebase emulator with paramater: `--export-on-exit`
- Add user, change options...
- Exit the emulator

> NB: The data will be saved in the `./firebase-local-data` folder. If you want to share the data, you can copy it to `./firebase-local-data.example` and commit it.

### How to reset the DB ?

`docker volume rm marble-backend_postgres-db` deletes the PG volume, useful to reset the app to a known state

In practice, this single-line will delete the stack and create a new one:
`docker compose down && docker volume rm marble-backend_postgres-db && docker compose up -d`
