package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// MyEvent represents input to the Lambda function.
type MyEvent struct {
	Name string `json:"name"`
}

// HandleRequest returns a pointer to a string containing invoke result
// and maybe an error.
func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	// otellambda does not have `WithMeterProvider` like with `WithTracerProvider`
	// https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda@v0.50.0#section-readme
	httpMetricExporter, _ := otlpmetrichttp.New(ctx)
	metricReader := sdkmetric.NewPeriodicReader(httpMetricExporter,
		// If internal is less than or equal to zero, 60s default is used
		// https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric#WithInterval
		sdkmetric.WithInterval(1*time.Second))
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(metricReader))
	// Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used, which fails to generate data.
	otel.SetMeterProvider(meterProvider)

	var meter = otel.Meter("foo-meter")
	lambdaStartCounter, _ := meter.Int64Counter("foo.tammy.test.count")
	lambdaStartCounter.Add(
		ctx,
		3)
	// Required with PeriodicReader interval 1s
	metricReader.ForceFlush(ctx)

	// Manually create child span
	_, span := tracer.Start(ctx, "my-manual-span")
	defer span.End()

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

	ctx := context.Background()
	httpTraceExporter, _ := otlptracehttp.New(ctx)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(httpTraceExporter))

	// set for manual trace creation
	tracer = traceProvider.Tracer("MyLambdaService")

	lambda.Start(
		otellambda.InstrumentHandler(
			HandleRequest,
			otellambda.WithTracerProvider(traceProvider)))
}
