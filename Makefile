run:
	 @npx tailwindcss -o ./src/styles/out.css
	 @templ generate
	 @go build -o ./tmp/main ./src/main.go
