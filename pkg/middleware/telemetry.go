package middleware

import (
	"clinicapp/pkg/logging"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func HandleRequestTelemetry(t *logging.Telemetry) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		startTime := time.Now().UTC()
		ctx.Next()
		duration := time.Since(startTime)

		status := strconv.Itoa(ctx.Writer.Status())

		request := appinsights.NewRequestTelemetry(ctx.Request.Method, ctx.Request.URL.Path, duration, status)

		request.Timestamp = time.Now().UTC()

		(*t.Client).Track(request)

	})
}

// func (t Telemetry) HandleRequestEventTelemetry(h func(*gin.Context)) gin.HandlerFunc {
// 	return gin.HandlerFunc(func(ctx *gin.Context) {

// 		// appinsights.EventTelemetry()

// 		startTime := time.Now().UTC()
// 		h(ctx)
// 		duration := time.Since(startTime)

// 		status := strconv.Itoa(ctx.Writer.Status())

// 		request := appinsights.NewRequestTelemetry(ctx.Request.Method, ctx.Request.URL.Path, duration, status)
// 		request := appinsights.NewEventTelemetry("LOGIN")
// 		request.Name = "LOGIN"

// 		request.Timestamp = time.Now().UTC()

// 		(*t.Client).Track(request)
// 		// (*t.Client).TrackEvent(request)
// 	})
// }

func trackError(t *logging.Telemetry, err error) {
	if err != nil {
		trace := appinsights.NewTraceTelemetry(err.Error(), appinsights.Error)
		trace.Timestamp = time.Now().UTC()
		(*t.Client).Track(trace)
	}
}

func trackWarning(t *logging.Telemetry, err error) {
	if err != nil {
		trace := appinsights.NewTraceTelemetry(err.Error(), appinsights.Warning)
		trace.Timestamp = time.Now().UTC()
		(*t.Client).Track(trace)
	}
}
