package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func handleRequestWithLog(h func(*gin.Context), telemetryClient appinsights.TelemetryClient) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		startTime := time.Now().UTC()
		h(ctx)
		duration := time.Since(startTime)

		status := strconv.Itoa(ctx.Writer.Status())

		request := appinsights.NewRequestTelemetry(ctx.Request.Method, ctx.Request.URL.Path, duration, status)

		request.Timestamp = time.Now().UTC()

		telemetryClient.Track(request)

	})
}
