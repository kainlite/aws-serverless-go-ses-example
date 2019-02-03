package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Response struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type Request struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

// This could be env vars
const (
	Sender    = "web@serverless.techsquad.rocks"
	Recipient = "kainlite@gmail.com"
	CharSet   = "UTF-8"
)

func lambdaHandler(ctx context.Context, req Request) (Response, error) {
	fmt.Printf("%+v\n", ctx)
	fmt.Printf("%+v\n", req)

	if len(req.Email) > 0 && len(req.Message) > 0 {
		send(req)
	}

	return Response{
		Id:      req.Email,
		Message: "Mail sent!",
	}, nil
}

func send(req Request) {
	// This could be an env var
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(req.Message),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(req.Message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(req.Email),
			},
		},
		// We are using the same sender because it needs to be validated in SES.
		Source: aws.String(Recipient),

		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)
}

func main() {
	lambda.Start(lambdaHandler)
}
