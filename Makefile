EXECUTABLE_NAME ?= server

generate:
	@echo "Generate everything"

build: migration_up generate
	@echo "Building server"
	@mkdir -p .bin
	@go build -o ./${EXECUTABLE_NAME} ./cmd/server

clean:
	@echo "deleting binaries"
	@rm ./server

migration_up:
	@echo "migration up"
	@go tool goose -env=.env -dir=migrations/ up

migration_down:
	@echo "migration down"
	@go tool goose -env=.env -dir=migrations/ down
