.PHONY: build deploy-heroku

build:
	GOOS=linux GOARCH=amd64 go build -o bin/bootstrap ./cmd/main.go
	echo "✅ Build successful."

deploy-heroku: build
	git add .
	git commit -m "Deploying latest version"
	git push heroku main
	echo "🚀 Deployment to Heroku completed."
