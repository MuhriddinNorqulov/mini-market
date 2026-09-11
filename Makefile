generate-docs:
	@swag init -g ./cmd/http/main.go -o ./src/infrastructure/docs

wire-build:
	@wiregenx --root ./src --out ../cmd/wire/provider.go
	@mkdir -p cmd/container
	@cd cmd/wire/ && wire && mv wire_gen.go ../container/container.go

docker-run.local:
	@cd docker && docker compose -f docker-compose.local.yml --env-file ./../env/.env up -d

