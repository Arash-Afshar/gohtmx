setup:
	npm install
	npm install -D tailwindcss
	go install github.com/air-verse/air@latest

dev:
	air -c .air.toml

dev-css:
	npx tailwindcss -i assets/styles.css -o static/styles.css --postcss --watch
