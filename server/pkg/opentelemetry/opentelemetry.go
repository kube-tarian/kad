package opentelemetry

import (
	"context"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tre "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

type Configurations struct {
	ServiceName  string `envconfig:"APPLICATION_NAME" default:"Applications"`
	CollectorURL string `envconfig:"OPTE_URL" default:"otelcollector.awsagent.optimizor.app:80"`
	IsInsecure   bool   `envconfig:"IS_OPTEL_INSECURE" default:"true"`
	IsEnabled    bool   `envconfig:"IS_OPTEL_ENABLED" default:"true"`
}

func getConfigurations() (opteConfig *Configurations, err error) {
	opteConfig = &Configurations{}
	if err = envconfig.Process("", opteConfig); err != nil {
		return nil, errors.WithStack(err)
	}
	return

}

func InitTracer() (func(context.Context) error, error) {
	config, err := getConfigurations()
	if err != nil {
		logger.Errorf("Unable to read open telemetry configurations")
		return nil, err
	}
	if !config.IsEnabled {
		return nil, nil
	}

	headers := map[string]string{
		"signoz-service-name": config.ServiceName,
	}
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if config.IsInsecure {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(config.CollectorURL),
			otlptracegrpc.WithHeaders(headers),
		),
	)

	if err != nil {
		logger.Errorf("unble initialize new object , error : %v", err)
		return nil, err
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", config.ServiceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.Errorf("Could not set resources:%v ", err)
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exporter.Shutdown, nil
}

var tracer tre.Tracer

func GetTracer(defaultValue string) tre.Tracer {
	if tracer == nil {
		serviceName := os.Getenv("APPLICATION_NAME")
		if serviceName == "" || serviceName == "Applications" {
			serviceName = defaultValue
		}
		tracer = otel.Tracer(serviceName)
	}
	return tracer
}

func BuildContext(ctx context.Context) context.Context {
	newCtx, _ := context.WithCancel(ctx)
	return newCtx
}
