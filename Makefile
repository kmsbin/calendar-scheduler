createDb:
	docker compose up -d --remove-orphans

upmigrate:
	migrate -path migrations/ -database "postgresql://kauli:1234@localhost:5432/calendator?sslmode=disable" -verbose up

destroyDb:
	docker compose down --remove-orphans