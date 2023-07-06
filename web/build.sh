uuid=$(uuidgen); GOOS=linux GOARCH=amd64 go build -ldflags="-X main.buildTag=$uuid"
