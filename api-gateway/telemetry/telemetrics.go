package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// global provider
var Meter metric.Meter

var MeterProvider *sdkmetric.MeterProvider

var Resource *resource.Resource

// resource provider from official docs adding.
func InitTelemetry() {

	Resource, err := resource.New(
		context.Background(),

		resource.WithTelemetrySDK(), //otel sdk version

		resource.WithProcess(), //adds process executable, runtime info

		resource.WithOS(),

		resource.WithHost(), //adds hostname and machine info

		resource.WithAttributes(
			semconv.ServiceName(
				"aegisflow",
			),

			semconv.ServiceVersion(
				"1.0.0",
			),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	exporter, err := prometheus.New()

	if err != nil {
		log.Fatal(err)
	}

	MeterProvider = //init telemetry
		sdkmetric.NewMeterProvider( //creates aegis flow metric engine
			sdkmetric.WithReader(
				exporter,
			),
			sdkmetric.WithResource(
				Resource,
			),
		)

	otel.SetMeterProvider( //register globally, every package can call
		MeterProvider, //otel.Meter and get the same provider
	)
	Meter = otel.Meter( //create meter
		"aegisflow",
	)
	Tracer = otel.Tracer( //initialise tracer
		"aegisflow",
	)

}
