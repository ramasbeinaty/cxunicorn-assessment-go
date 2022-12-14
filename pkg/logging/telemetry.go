package logging

import (
	"log"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func NewTelemetry() *Telemetry {
	// define an azure app insight telemetry client

	t := new(Telemetry)

	_client := appinsights.NewTelemetryClient(os.Getenv("INSTRUMENTATION_KEY"))

	t.Client = &_client

	/*Set role instance name globally -- this is usually the name of the service submitting the telemetry*/
	// t.Client.Context().Tags.Cloud().SetRole("clinic_app")
	(*t.Client).Context().Tags.Cloud().SetRole("clinic_app")

	/*turn on diagnostics to help troubleshoot problems with telemetry submission. */
	appinsights.NewDiagnosticsMessageListener(func(msg string) error {
		log.Printf("[%s] %s\n", time.Now().Format(time.UnixDate), msg)
		return nil
	})

	return t
}

type Telemetry struct {
	Client *appinsights.TelemetryClient
}
