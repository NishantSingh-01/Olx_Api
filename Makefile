.PHONY: build run start clean

build:
	@go build -o bin/api/main.exe ./cmd/api/main.go

run:
	@go run ./cmd/api/main.go

start:
	@.\bin\api\main.exe

clean:
	@if exist bin rmdir /s /q bin