package main

import (
    "fmt"
    "context"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ses"
    "github.com/aws/aws-sdk-go/aws/awserr"
)

type Request struct {
    Email string `json:"email"`
    Message string `json:"message"`
}

type Response struct {
    Message string `json:"message"`

}

type Event struct {
    Email string `json:"email"`
    Message string `json:"message"`
}

// This could be env vars
const (
    Sender = "web@skynetng.pw"
    Recipient = "kainlite@gmail.com"
    Subject = "New mail from the site..."
    CharSet = "UTF-8"
)

func Handler(ctx context.Context, event Event) (Response, error) {
        // To get some debug info
        // fmt.Printf("%+v\n", ctx)
        // fmt.Printf("%+v\n", ev)

	send(event)

	return Response{
		Message: "Mail sent!",
	}, nil
}

func send(event Event) {
    // This could be an env var
    sess, err := session.NewSession(&aws.Config{
        Region:aws.String("us-east-1")},
    )

    // Create an SES session.
    svc := ses.New(sess)

    // Assemble the email.
    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            CcAddresses: []*string{
            },
            ToAddresses: []*string{
                aws.String(Recipient),
            },
        },
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(event.Message),
                },
                Text: &ses.Content{
                    Charset: aws.String(CharSet),
                    Data:    aws.String(event.Message),
                },
            },
            Subject: &ses.Content{
                Charset: aws.String(CharSet),
                Data:    aws.String(Subject),
            },
        },
        Source: aws.String(event.Email),
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
	lambda.Start(Handler)
}
