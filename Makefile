run:
	 @npx tailwindcss -o ./public/styles/out.css
	 @templ generate
	 @go build -o ./tmp/main ./src/main.go
format:
	@gofmt -w .
