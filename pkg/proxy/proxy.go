package proxy

import (
	"log/slog"

	"github.com/josephgoksu/gok-proxy/pkg/handler"

	"github.com/josephgoksu/gok-proxy/pkg/config"

	"github.com/valyala/fasthttp"
)

type ProxyServer struct {
	cfg    *config.Config
	logger *slog.Logger
	server *fasthttp.Server
}

func NewProxyServer(cfg *config.Config, logger *slog.Logger) *ProxyServer {
	requestHandler := handler.NewRequestHandler(logger)

	return &ProxyServer{
		cfg:    cfg,
		logger: logger,
		server: &fasthttp.Server{
			Handler:            requestHandler.HandleRequest,
			MaxConnsPerIP:      cfg.MaxConnections,
			MaxRequestsPerConn: cfg.MaxRequestsPerConn,
		},
	}
}

func (p *ProxyServer) Start() error {
	return p.server.ListenAndServe(p.cfg.ServerAddress)
}

func (p *ProxyServer) Shutdown() error {
	return p.server.Shutdown()
}
