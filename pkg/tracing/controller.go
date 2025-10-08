package tracing

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

const (
	enableExporterEnv    = "ENABLE_OTEL_EXPORTER"
	collectorEndpointEnv = "OTEL_COLLECTOR_ENDPOINT"
	insecureExporterEnv  = "OTEL_EXPORTER_OTLP_INSECURE"
	defaultServiceName   = "kserve-controller"
	tracerSetupTimeout   = 10 * time.Second
)

var noopShutdown = func(context.Context) error { return nil }

// SetupControllerTracing configures the global OpenTelemetry tracer provider for the
// controller if exporter support is enabled.
func SetupControllerTracing(ctx context.Context, log logr.Logger) (func(context.Context) error, error) {
	if !isTruthy(os.Getenv(enableExporterEnv)) {
		log.Info("OpenTelemetry exporter disabled; controller spans will not be exported", "env", enableExporterEnv)
		return noopShutdown, nil
	}

	exporter, endpoint, insecure, err := buildExporter(ctx)
	if err != nil {
		return noopShutdown, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithAttributes(semconv.ServiceName(defaultServiceName)),
	)
	if err != nil {
		return noopShutdown, fmt.Errorf("failed to construct resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(provider)
	log.Info("Configured OpenTelemetry tracing for controller spans", "endpoint", endpoint, "insecure", insecure)

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, tracerSetupTimeout)
		defer cancel()

		if err := provider.Shutdown(ctx); err != nil {
			log.Error(err, "failed to shut down tracer provider")
			return err
		}
		return nil
	}, nil
}

func buildExporter(parent context.Context) (sdktrace.SpanExporter, string, bool, error) {
	ctx, cancel := context.WithTimeout(parent, tracerSetupTimeout)
	defer cancel()

	endpoint := strings.TrimSpace(os.Getenv(collectorEndpointEnv))
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	}

	dialOpts := []otlptracegrpc.Option{}

	if endpoint != "" {
		dialOpts = append(dialOpts, otlptracegrpc.WithEndpoint(endpoint))
	}

	insecure := isTruthyWithDefault(os.Getenv(insecureExporterEnv), true)
	if insecure {
		dialOpts = append(dialOpts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, dialOpts...)
	if err != nil {
		return nil, endpoint, insecure, err
	}

	return exporter, endpoint, insecure, nil
}

func isTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func isTruthyWithDefault(value string, fallback bool) bool {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return isTruthy(value)
}
