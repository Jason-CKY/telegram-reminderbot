REMINDERBOT_VERSION ?= 1.16.0
BACKUP_DIR ?= ~/backup
REMINDERBOT_BACKUP_FILE ?= reminderbot-backup.tar

.DEFAULT_GOAL := help

# declares .PHONY which will run the make command even if a file of the same name exists
.PHONY: help
help:			## Help command
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint:			## Lint check
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src odinuge/yapf:latest yapf \
	--style '{based_on_style: pep8, dedent_closing_brackets: true, coalesce_brackets: true}' \
	--no-local-style --verbose --recursive --diff --parallel app compose

format:			## Format code in place to conform to lint check
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src odinuge/yapf:latest yapf \
	--style '{based_on_style: pep8, dedent_closing_brackets: true, coalesce_brackets: true}' \
	--no-local-style --verbose --recursive --in-place --parallel app compose

pyflakes:		## Pyflakes check for any unused variables/classes
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src python:3.8 \
	/bin/bash -c "pip install --upgrade pyflakes && python -m pyflakes /src && echo 'pyflakes passed!'"

backup-reminderbot:		## Backup database volumes to BACKUP_DIR
	docker stop reminderbot reminderbot_db
	docker run --rm --volumes-from reminderbot_db -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /data/db && tar cvf /backup/$(REMINDERBOT_BACKUP_FILE) ."
	docker start reminderbot reminderbot_db

restore-backup-reminderbot:		## Restore volumes backup from BACKUP_DIR
	ls $(BACKUP_DIR) | grep reminderbot-backup.tar
	docker volume create telegram-reminderbot_reminderbot-db
	docker run --rm -v telegram-reminderbot_reminderbot-db:/recover -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /recover && tar xvf /backup/reminderbot-backup.tar"

start-prod:			## Pull latest version of reminderbot image and run docker-compose up
	docker-compose pull reminderbot
	docker-compose up -d
	
start-dev:			## Run dev instance of reminderbot with live reload of api
	docker-compose -f docker-compose.dev.yml up --build -d

build-reminderbot:		## Build docker image for reminderbot and pollingserver
	docker build --tag jasoncky96/telegram-reminderbot:latest -f ./compose/webserver/Dockerfile .
	docker build --tag jasoncky96/telegram-pollingserver:latest -f ./compose/pollingserver/Dockerfile .

deploy-pollingserver:	## Deploy docker image for pollingserver
	docker buildx build --push --tag jasoncky96/telegram-pollingserver:latest --file ./compose/pollingserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

deploy-reminderbot:		## Deploy docker image for reminderbot
	docker buildx build --push --tag jasoncky96/telegram-reminderbot:$(REMINDERBOT_VERSION) --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .
	docker buildx build --push --tag jasoncky96/telegram-reminderbot:latest --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

stop-reminderbot:		## Run docker-compose down
	docker-compose down
	
destroy:				## Run docker-compose down -v
	docker-compose down -v
