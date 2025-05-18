package main

import (
	stdlog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	applogger "github.com/josephgoksu/gok-proxy/pkg/log"
	"github.com/josephgoksu/gok-proxy/pkg/pool"
	"github.com/josephgoksu/gok-proxy/pkg/proxy"

	"github.com/josephgoksu/gok-proxy/pkg/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		stdlog.Fatalf("Error loading config: %v", err)
	}

	// Setup logging
	logger, err := applogger.NewLogger(cfg.LogLevel)
	if err != nil {
		stdlog.Fatalf("Error setting up logger: %v", err)
	}
	slog.SetDefault(logger) // Set the global logger

	// Initialize connection pool with config
	pool.InitConnPool(cfg)

	// Setup metrics
	http.Handle("/metrics", promhttp.Handler())

	// Initialize and start the proxy server
	proxyServer := proxy.NewProxyServer(cfg, logger) // Pass the slog.Logger directly
	logger.Info("Starting proxy server", "address", cfg.ServerAddress)

	// Channel to handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := proxyServer.Start(); err != nil {
			logger.Error("Failed to start proxy server", "error", err)
			os.Exit(1) // Exit if server fails to start
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	<-quit
	logger.Info("Shutting down proxy server...")
	if err := proxyServer.Shutdown(); err != nil {
		logger.Error("Failed to gracefully shut down proxy server", "error", err)
		os.Exit(1) // Exit if shutdown fails
	}
	logger.Info("Proxy server stopped")
}
