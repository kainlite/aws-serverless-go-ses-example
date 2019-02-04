.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o send_mail/send_mail send_mail/main.go

clean:
	rm -rf ./send_mail/send_mail

deploy: clean build
	sls deploy --verbose
