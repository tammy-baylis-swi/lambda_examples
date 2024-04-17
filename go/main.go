package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

// MyEvent represents input to the Lambda function.
type MyEvent struct {
	Name string `json:"name"`
}

// HandleRequest returns a pointer to a string containing invoke result
// and maybe an error.
func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received a nil event")
	}
	message := fmt.Sprintf("Hello %s!", event.Name)
	return &message, nil
}

// main is the entry point to our Lambda function.
func main() {
	log.Println("Starting HandleRequest")
	lambda.Start(HandleRequest)
}
