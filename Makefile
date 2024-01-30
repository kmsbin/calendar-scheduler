
createDb:
	docker compose up -d --remove-orphans

upmigrate:
	migrate -path src/database/migrations/ -database "postgresql://root:1234@localhost:5432/skill_share?sslmode=disable" -verbose up

destroyDb:
	docker compose down --remove-orphans