build:
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-sdk-go
	env GOOS=linux go build -ldflags="-s -w" -o bin/send_mail send_mail/main.go
