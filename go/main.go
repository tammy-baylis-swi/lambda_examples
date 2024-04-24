package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
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
	startTime := time.Now().UnixMilli()

	// Adding a 500ms sleep for histogram testing
	time.Sleep(1 * time.Millisecond)

	// otellambda does not have `WithMeterProvider` like with `WithTracerProvider`
	// so we're setting up MeterProvider here inside the handler function
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

	// Test counter and histogram
	var meter = otel.Meter("foo-meter")
	testCounter, _ := meter.Int64Counter("foo.tammy.test.count")
	testCounter.Add(
		ctx,
		3)
	testHistogram, _ := meter.Int64Histogram(
		"foo.tammy.test.histo",
		metric.WithDescription("The duration of handler execution."),
		metric.WithUnit("ms"))

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

	endTime := time.Now().UnixMilli()
	testHistogram.Record(
		ctx,
		(endTime - startTime))

	// Required with PeriodicReader interval 1s
	metricReader.ForceFlush(ctx)

	// One more debug line
	recordTime := endTime - startTime
	timeMessage := fmt.Sprintf("endTime minus startTime was %s", strconv.Itoa(int(recordTime)))
	log.Println(timeMessage)

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
