package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
)

// MyEvent represents input to the Lambda function.
type MyEvent struct {
	Name string `json:"name"`
}

// HandleRequest returns a pointer to a string containing invoke result
// and maybe an error.
func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		message := "received a nil event"
		log.Println(message)
		return nil, fmt.Errorf(message)
	}
	message := fmt.Sprintf("Hello %s!", event.Name)
	log.Println(message)
	return &message, nil
}

// main is the entry point to our Lambda function,
// wrapped by OpenTelemetry Go SDK.
func main() {
	log.Println("Starting HandleRequest")

	httpTraceExporter, _ := otlptracehttp.New(context.Background())
	traceProvider := trace.NewTracerProvider(
		trace.WithSyncer(httpTraceExporter))

	lambda.Start(
		otellambda.InstrumentHandler(
			HandleRequest,
			otellambda.WithTracerProvider(traceProvider)))
}
