package test

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/josephgoksu/gok-proxy/pkg/config"
	"github.com/josephgoksu/gok-proxy/pkg/handler"
	"github.com/josephgoksu/gok-proxy/pkg/pool"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestHandleRequest_HTTP_Get_ExternalSite(t *testing.T) {
	// Initialize connection pool with a test configuration
	testConfig := &config.Config{
		ClientReadTimeoutSeconds:         10, // Sensible timeout for a test
		ClientWriteTimeoutSeconds:        10,
		ClientMaxIdleConnDurationSeconds: 30,
		// Other server-specific config fields are not strictly necessary for InitConnPool
		// but good to have for completeness if default values are not 0 or empty.
		ServerAddress:      ":0",    // Not used by client pool
		MaxConnections:     100,     // Not used by client pool
		LogLevel:           "debug", // Not used by client pool
		MaxRequestsPerConn: 100,     // Not used by client pool
	}
	pool.InitConnPool(testConfig)

	// Create a discard logger for testing
	discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	h := handler.NewRequestHandler(discardLogger)

	// Target URL for the proxy request (user changed to this)
	// This makes it an integration test, requiring network access.
	targetURL := "https://ifconfig.io"
	req := httptest.NewRequest("GET", targetURL, nil)

	fastReq := &fasthttp.Request{}
	fastReq.SetRequestURI(req.URL.String()) // fasthttp client needs the full URL for proxying this way

	fastCtx := &fasthttp.RequestCtx{}
	fastCtx.Init(fastReq, nil, nil) // Init with the request
	// The actual remote address is nil, which is fine for this type of test
	// where the handler doesn't rely on client IP for this path.

	h.HandleRequest(fastCtx)

	// Ensure the response status code is OK (200)
	// This implies ifconfig.io returned 200 through the proxy
	assert.Equal(t, fasthttp.StatusOK, fastCtx.Response.StatusCode(), "Expected status OK")

	// Optionally, you could check parts of the body if you know the expected content,
	// but for ifconfig.io, the IP address in the body changes, so a status check is primary.
	// e.g., assert.NotEmpty(t, fastCtx.Response.Body(), "Response body should not be empty")
}

// TODO: Consider adding a test for the CONNECT method (handleTunneling)
// This would require a different setup, likely involving a mock destination server.
