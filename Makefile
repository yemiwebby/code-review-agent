.PHONY: build deploy

build:
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap ./cmd/main.go
	echo "✅ Build successful."

deploy: build
	git add .
	git commit -m "Deploying latest version"
	git push heroku main