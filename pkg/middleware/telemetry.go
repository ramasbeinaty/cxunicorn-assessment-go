package middleware

import (
	"log"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func NewTelemetryClient() *telemetry {
	// define an azure app insight telemetry client

	t := new(telemetry)

	t.Client = appinsights.NewTelemetryClient(os.Getenv("INSTRUMENTATION_KEY"))

	/*Set role instance name globally -- this is usually the name of the service submitting the telemetry*/
	t.Client.Context().Tags.Cloud().SetRole("clinic_app")

	/*turn on diagnostics to help troubleshoot problems with telemetry submission. */
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})

	return t
}

type telemetry struct {
	Client appinsights.TelemetryClient
}
