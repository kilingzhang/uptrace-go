package spanexp

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

type Option func(*Config)

// SpanFilter is a function that is used to filter and change Uptrace spans.
type SpanFilter func(*Span) bool

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// `service.name` resource attribute. It is merged with Config.Resource.
	ServiceName string
	// `service.version` resource attribute. It is merged with Config.Resource.
	ServiceVersion string
	// Any other resource attributes. They are merged with Config.Resource.
	ResourceAttributes []label.KeyValue
	// Resource contains attributes representing an entity that produces telemetry.
	// These attributes are copied to all spans and events.
	//
	// The default is `resource.New`.
	Resource *resource.Resource

	// Filters are functions that are used to filter and change Uptrace spans.
	Filters []SpanFilter

	// Global TextMapPropagator used by OpenTelemetry.
	// The default is propagation.TraceContext and propagation.Baggage.
	TextMapPropagator propagation.TextMapPropagator

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client
	// Max number of retries when sending data to Uptrace.
	// The default is 3.
	MaxRetries int

	// Name of the tracer used by Uptrace client.
	// The default is github.com/uptrace/uptrace-go.
	TracerName string

	// PrettyPrint pretty prints spans to the stdout.
	PrettyPrint bool

	// When specified it overwrites the default Uptrace tracer provider.
	// It can be used to configure Uptrace client to use OTLP exporter.
	TracerProvider trace.TracerProvider

	// Disabled disables the exporter.
	// The default is to use UPTRACE_DISABLED environment var.
	Disabled bool

	// Trace enables Uptrace exporter instrumentation.
	Trace bool

	// ClientTrace enables httptrace instrumentation on the HTTP client used by Uptrace.
	ClientTrace bool

	inited bool
}

func (cfg *Config) Init(opts ...Option) {
	if cfg.inited {
		return
	}
	cfg.inited = true

	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		cfg.Disabled = true
		return
	}

	if cfg.DSN == "" {
		if dsn, ok := os.LookupEnv("UPTRACE_DSN"); ok {
			cfg.DSN = dsn
		}
	}

	if cfg.Resource == nil {
		resource, err := resource.New(context.TODO())
		if err == nil {
			cfg.Resource = resource
		}
	}

	{
		kvs := cfg.ResourceAttributes

		if cfg.ServiceName != "" {
			kvs = append(kvs, semconv.ServiceNameKey.String(cfg.ServiceName))
		}
		if cfg.ServiceVersion != "" {
			kvs = append(kvs, semconv.ServiceNameKey.String(cfg.ServiceName))
		}

		if len(kvs) > 0 {
			cfg.Resource = resource.Merge(
				resource.NewWithAttributes(kvs...),
				cfg.Resource,
			)
		}
	}

	if cfg.TextMapPropagator == nil {
		cfg.TextMapPropagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	if cfg.TracerName == "" {
		cfg.TracerName = "github.com/uptrace/uptrace-go"
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	switch cfg.MaxRetries {
	case -1:
		cfg.MaxRetries = 0
	case 0:
		cfg.MaxRetries = 3
	}

	if cfg.ClientTrace {
		cfg.Trace = true
	}

	for _, opt := range opts {
		opt(cfg)
	}
}
