//for traces from OpenTelemetry
//add a trace exporter
package telemetry

import (


	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	stdouttrace "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

)

var Tracer trace.Tracer

var TracerProvider *sdktrace.TracerProvider

func InitTracer() error {

	exporter, err :=
		stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
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
            Resource,
        ),
    )

	otel.SetTracerProvider(
		TracerProvider,
	)

	Tracer =
		TracerProvider.Tracer(
			"aegisflow",
		)

	return nil
}

