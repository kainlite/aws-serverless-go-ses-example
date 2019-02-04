package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Response events.APIGatewayProxyResponse

type RequestData struct {
	Email   string
	Message string
}

// This could be env vars
const (
	Sender    = "web@serverless.techsquad.rocks"
	Recipient = "kainlite@gmail.com"
	CharSet   = "UTF-8"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Printf("Request: %+v\n", request)

	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))

	var requestData RequestData
	json.Unmarshal([]byte(request.Body), &requestData)

	fmt.Printf("RequestData: %+v", requestData)
	var result string
	if len(requestData.Email) > 0 && len(requestData.Message) > 0 {
		result, _ = send(requestData.Email, requestData.Message)
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            result,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "send-mail-handler",
		},
	}

	return resp, nil
}

func send(Email string, Message string) (string, error) {
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
					Data:    aws.String(Message),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Email),
			},
		},
		// We are using the same sender because it needs to be validated in SES.
		Source: aws.String(Sender),

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

		return "there was an unexpected error", err
	}

	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)
	return "sent!", err
}

func main() {
	lambda.Start(Handler)
}
