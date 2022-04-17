REMINDERBOT_VERSION ?= 1.16
BACKUP_DIR ?= ~/backup
REMINDERBOT_BACKUP_FILE ?= reminderbot-backup.tar

format: 
	yapf -i -r -p .

backup-reminderbot:
	docker stop reminderbot reminderbot_db
	docker run --rm --volumes-from reminderbot_db -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /data/db && tar cvf /backup/$(REMINDERBOT_BACKUP_FILE) ."
	docker start reminderbot reminderbot_db

restore-backup-reminderbot:
	ls $(BACKUP_DIR) | grep reminderbot-backup.tar
	docker volume create telegram-reminderbot_reminderbot-db
	docker run --rm -v telegram-reminderbot_reminderbot-db:/recover -v $(BACKUP_DIR):/backup ubuntu bash -c "cd /recover && tar xvf /backup/reminderbot-backup.tar"

start-prod:
	docker-compose pull reminderbot
	docker-compose up -d
	
start-dev:
	docker-compose -f docker-compose.dev.yml up --build -d

build-reminderbot:
	docker build --tag jasoncky96/telegram-reminderbot:latest -f ./compose/webserver/Dockerfile .
	docker build --tag jasoncky96/telegram-pollingserver:latest -f ./compose/pollingserver/Dockerfile .

deploy-pollingserver:
	docker buildx build --push --tag jasoncky96/telegram-pollingserver:latest --file ./compose/pollingserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

deploy-reminderbot:
	docker buildx build --push --tag jasoncky96/telegram-reminderbot:$(REMINDERBOT_VERSION) --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .
	docker buildx build --push --tag jasoncky96/telegram-reminderbot:latest --file ./compose/webserver/Dockerfile --platform linux/arm/v7,linux/arm64/v8,linux/amd64 .

stop-reminderbot:
	docker-compose down
	
destroy:
	docker-compose down -v
