# Go htmx server

## Dependencies

- docker (docker-desktop if you are using windows)
- docker-compose (comes with docker-desktop, but can install [here](https://docs.docker.com/compose/install/standalone/) if you are not on windows)
- [>=Go v1.21](https://go.dev/doc/install)
- [Air](https://github.com/cosmtrek/air)

## Features

- [air](https://github.com/cosmtrek/air) for code reloading in dev environment
- [Directus](https://directus.io/) for headless CMS and API routes for CRUD operations

## Quickstart (development mode)

Run `cp .env.example .env`, and fill in the relevant information

```sh
make build-dev
# start golang server with code reloading using air
air
```

### Todos

- `Reminder in Construction` and `Reminder` collections in directus
- Delete reminder in construction when cancelled
- Finish the flow for making once-off reminder
- Finish flow for recurring reminders
- Handle repeated `/remind` cases (ignore them)
- Handle group chat settings
  - Handle reminder from each user separately
