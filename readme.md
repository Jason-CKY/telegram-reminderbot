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

- [X] Handle repeated `/remind` cases (delete old reminder in construction and create new one)
- [X] Finish flow for daily reminders
- [X] Finish flow for weekly reminders
- [X] Finish flow for monthly reminders
- [X] Finish the flow for making once-off reminder
- [X] Finish flow for yearly reminders
- [X] Set separate go-routine to check if any reminder is due
- [X] Option to renew reminders when they are triggered
- [X] Handle chat settings  (require separate table to store settings)
- [X] Handle listing and deleting reminders
- [X] Handle image reminders
- [ ] Handle group reminders
- [ ] Handle group settings
