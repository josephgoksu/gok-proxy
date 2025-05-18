package handler

import (
	"io"
	"log/slog"
	"net"

	"github.com/josephgoksu/gok-proxy/pkg/metrics"
	"github.com/josephgoksu/gok-proxy/pkg/pool"

	"github.com/valyala/fasthttp"
)

// RequestHandler holds dependencies for request handling, like a logger.
type RequestHandler struct {
	logger *slog.Logger
}

// NewRequestHandler creates a new RequestHandler with the given logger.
func NewRequestHandler(logger *slog.Logger) *RequestHandler {
	return &RequestHandler{logger: logger}
}

// HandleRequest handles both normal HTTP requests and HTTP CONNECT requests for HTTPS
func (h *RequestHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	metrics.IncrementRequestCounter()

	// Handle HTTP CONNECT method for HTTPS proxying
	if string(ctx.Method()) == fasthttp.MethodConnect {
		h.handleTunneling(ctx)
	} else {
		h.handleHTTP(ctx)
	}
}

// handleHTTP handles standard HTTP requests
func (h *RequestHandler) handleHTTP(ctx *fasthttp.RequestCtx) {
	client := pool.GetConnection()
	defer pool.PutConnection(client)

	req := &ctx.Request
	resp := &ctx.Response

	if err := client.Do(req, resp); err != nil {
		h.logger.Error("Failed to process HTTP request", "error", err, "url", string(req.URI().FullURI()))
		ctx.Error("Failed to process request: "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}
}

// handleTunneling handles HTTP CONNECT requests
func (h *RequestHandler) handleTunneling(ctx *fasthttp.RequestCtx) {
	destinationHost := string(ctx.Host())
	destinationConn, err := net.Dial("tcp", destinationHost)
	if err != nil {
		h.logger.Error("Failed to connect to destination for CONNECT", "error", err, "host", destinationHost)
		ctx.Error("Failed to connect to destination: "+err.Error(), fasthttp.StatusServiceUnavailable)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBodyRaw(nil) // No body for 200 OK on CONNECT

	ctx.Hijack(func(clientConn net.Conn) {
		defer clientConn.Close()
		defer destinationConn.Close()

		// Use a channel to wait for both copy operations to finish
		done := make(chan struct{}, 2)

		go func() {
			_, copyErr := io.Copy(destinationConn, clientConn)
			if copyErr != nil {
				h.logger.Warn("Error copying from client to destination", "error", copyErr, "host", destinationHost)
			}
			done <- struct{}{}
		}()

		go func() {
			_, copyErr := io.Copy(clientConn, destinationConn)
			if copyErr != nil {
				h.logger.Warn("Error copying from destination to client", "error", copyErr, "host", destinationHost)
			}
			done <- struct{}{}
		}()

		<-done
		<-done
	})
}
