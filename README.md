# OTEL GO instrumentation with GO ECHO

example MVP of instrumenting Open Telemetry in GO

## In microservice environment make http call to another serivce like this

> Disclaimer: this once worked, need to verify properly

```go
import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func callServiceA(ctx context.Context) error {
	url := "http://service-a:8080/endpoint"

	// initialize client
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// create request
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	// Inject trace context into the outgoing request
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
```

### Deployment Notes

the following env variable must be provided
```sh
RUNTIME_ENV
SERVICE_NAME
OTLP_ENDPOINT
```

