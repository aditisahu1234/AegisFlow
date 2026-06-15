//for traces from OpenTelemetry
//add a trace exporter
package telemetry

import (


	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	stdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

)

var Tracer trace.Tracer

var TracerProvider *sdktrace.TracerProvider

func InitTracer() error {

	exporter, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)

	res := resource.NewWithAttributes(	//deployement environment
		"",
		attribute.String(
			"service.name",
			"aegisflow",
		),
	
		attribute.String(
			"deployment.environment",
			"development",
		),
		attribute.String(
			"service.version",
			"1.0.0",
		),
	)

	if err != nil {
		return err
	}

	TracerProvider =
    sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(
            exporter,
        ),

        sdktrace.WithResource(
            res,
        ),
    )

	otel.SetTracerProvider(
		TracerProvider,
	)

	otel.SetTextMapPropagator(	//trace context propagation
		propagation.TraceContext{},
	)

	Tracer =
		TracerProvider.Tracer(
			"aegisflow",
		)

	return nil
}

