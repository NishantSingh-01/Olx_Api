.PHONY: build run start clean migrate-up migrate-down

build:
	@go build -o bin/api/main.exe ./cmd/api/main.go

run: build
	@.\bin\api\main.exe
	

start:
	@.\bin\api\main.exe


migrate-up:
	@go run ./cmd/migrate/main.go up

migrate-down:
	@go run ./cmd/migrate/main.go down

migrate-force:
	@go run ./cmd/migrate/main.go force $(version)

migrate-drop:
	@go run ./cmd/migrate/main.go drop


clean:
	@if exist bin rmdir /s /q bin