package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	// Your app handler.
	var handler http.Handler
	handler = http.HandlerFunc(userProfileEndpoint)

	// Wrap it with OpenTelemetry plugin.
	handler = otelhttp.WithRouteTag("/profiles/:username", handler)
	handler = otelhttp.NewHandler(handler, "server-name")

	// Register handler.
	http.Handle("/profiles/", handler)

	srv := &http.Server{
		Addr:    ":9999",
		Handler: handler,
	}

	fmt.Println("listening on http://localhost:9999")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func userProfileEndpoint(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	username := strings.Replace(req.URL.Path, "/profiles/", "", 1)

	name, err := selectUser(ctx, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, err.Error())
		return
	}

	fmt.Fprintf(w, `<html><h1>Hello %s %s </h1></html>`+"\n", username, name)
}

func selectUser(ctx context.Context, username string) (string, error) {
	_, span := tracer.Start(ctx, "selectUser")
	defer span.End()

	span.SetAttributes(attribute.String("username", username))

	if username == "admin" {
		return "Joe", nil
	}

	return "", fmt.Errorf("username=%s not found", username)
}
