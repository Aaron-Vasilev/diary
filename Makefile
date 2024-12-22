run:
	@templ generate
	@go build -o ./tmp/main ./src/main.go
build:
	@go build -o ./tmp/main ./src/main.go
dev:
	@npx concurrently "air" "npx tailwindcss -o ./public/styles/out.css --watch"
format:
	@gofmt -w .
	@templ fmt .
start:
	@supervisord -c ./supervisord.conf
launch:
	@templ generate
	@go build -o ./tmp/main ./src/main.go
	@echo Build ends
	@sudo supervisorctl restart diary
	@echo Started
restart:
	@sudo supervisorctl restart diary
stop:
	@supervisorctl shutdown
