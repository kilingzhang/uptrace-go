package main

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/uptrace/uptrace-go/extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
		uptrace.WithPrettyPrintSpanExporter(),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)


	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger = otelzap.Wrap(logger, otelzap.WithLevel(zap.NewAtomicLevelAt(zap.ErrorLevel)))

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")

	// You must use Ctx to propagate the active span.
	logger.Ctx(ctx).Error("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))

	span.End()

	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
