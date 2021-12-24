package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/komem3/logplug"
	"github.com/komem3/logplug/gcpopt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	port      = "9090"
	projectID = "test-project"
)

func init() {
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}
	if metadata.OnGCE() {
		p, err := metadata.ProjectID()
		if err != nil {
			log.Fatalf("[CRITICAL] get project: %v", err)
		}
		projectID = p
	}
}

func NewLog(ctx context.Context) *log.Logger {
	l := log.New(log.Writer(), log.Prefix(), log.Flags())
	if metadata.OnGCE() {
		span := trace.SpanFromContext(ctx).SpanContext()
		l.SetPrefix(fmt.Sprintf("[logging.googleapis.com/trace:projects/%s/traces/%s][logging.googleapis.com/spanId:%s][logging.googleapis.com/trace_sampled:%t]%s",
			projectID, span.TraceID(), span.SpanID(), span.IsSampled(), l.Prefix()))
	}
	return l
}

func main() {
	if metadata.OnGCE() {
		log.SetOutput(logplug.NewJSONPlug(os.Stderr, gcpopt.NewGCPOptions("DEBUG")...))
		log.SetFlags(gcpopt.LogFlags)
	}

	ctx := context.Background()

	{
		exporter, err := cloudtrace.New(cloudtrace.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("[CRITICAL] new cloudtrace")
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
		)

		otel.SetTracerProvider(tp)

		defer tp.Shutdown(ctx)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l := NewLog(r.Context())

		l.Printf("[DBG] request serve %s", r.RequestURI)
		l.Print("default level is info")
		l.Print("[ERR] error log")
		l.Print("[ALERT] alart")

		w.Write([]byte("hello world"))
	})

	otelHandler := otelhttp.NewHandler(mux, "Hello", otelhttp.WithPropagators(propagator.New()))

	log.Printf("[DBG] listen %s", port)
	log.Panicf("[CRITICAL] %v", http.ListenAndServe(":"+port, otelHandler))
}
