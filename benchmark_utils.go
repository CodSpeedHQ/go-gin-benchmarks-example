package api

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// DummyResponseWriter implements http.ResponseWriter but discards all data
// This eliminates overhead from httptest.NewRecorder() in benchmarks
type DummyResponseWriter struct{}

func (d *DummyResponseWriter) Header() http.Header {
	return http.Header{}
}

func (d *DummyResponseWriter) Write(data []byte) (int, error) {
	// Discard all data - do nothing
	return len(data), nil
}

func (d *DummyResponseWriter) WriteHeader(statusCode int) {
	// Do nothing - discard status code
}

// setupBenchmarkRouter wraps the main setupRouter with benchmark mode configuration
func setupBenchmarkRouter() *gin.Engine {
	// Set Gin to release mode for benchmarks
	gin.SetMode(gin.ReleaseMode)
	// Discard all output during benchmarks to only preserve benchmark output
	gin.DefaultWriter = io.Discard
	return setupRouter()
}

func benchmarkRequest(b *testing.B, req *http.Request) {
	router := setupBenchmarkRouter()
	w := new(DummyResponseWriter)
	for b.Loop() {
		router.ServeHTTP(w, req)
	}
}
