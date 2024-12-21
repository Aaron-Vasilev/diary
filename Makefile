run:
	 @templ generate
	 @go build -o ./tmp/main ./src/main.go
dev:
	@npx concurrently "air" "npx tailwindcss -o ./public/styles/out.css --watch"
format:
	@gofmt -w .
	@templ fmt .
start:
	@supervisord -c ./supervisord.conf
launch:
	@sudo supervisorctl shutdown
	@go build -o ./tmp/bot ./main.go
	@echo Build ends
	@sudo supervisord -c ./supervisord.conf
	@echo Started
stop:
	@supervisorctl shutdown
