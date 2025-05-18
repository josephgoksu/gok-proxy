<img src="./assets/logo-text.png" alt="gok-proxy" height="250" width="250" align="right" style="margin-left: 10px;" />

# gok-proxy: Lightweight High-Performance Go Proxy Server

`gok-proxy` (Gök is "sky" in Turkish) is a lightweight, exceptionally performant, and scalable HTTP/HTTPS proxy server built with Go. It leverages `fasthttp` for speed and efficiency, offering a modern alternative to traditional proxy solutions like Squid.

## Core Features

- **Exceptional Performance**: Built on Go's concurrency model and the high-speed `fasthttp` library for optimal resource utilization and request handling.
- **HTTP/1.1 Support**: Full support for HTTP/1.1 protocols.
- **HTTPS Proxying (CONNECT Tunneling)**: Securely proxies HTTPS traffic using the HTTP CONNECT method.
- **Structured Logging**: Employs `slog` (Go's standard structured logging package) for clear, configurable, and machine-parsable logs.
- **Prometheus Metrics**: Integrates seamlessly with Prometheus for comprehensive operational metrics and monitoring.
- **Flexible Configuration**: Utilizes `viper` for easy configuration management through a `config.yaml` file, with sensible defaults.
- **Graceful Shutdown**: Ensures clean termination by handling OS signals, preventing abrupt disconnections.
- **Connection Pooling**: Optimizes outgoing client connections using `sync.Pool` for `fasthttp.Client` instances.

## Prerequisites

- Go 1.24.2 or newer.
- A Unix-like environment (Linux, macOS) is recommended for development and deployment.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/josephgoksu/gok-proxy.git
cd gok-proxy
```

### 2. Install Dependencies

Go modules will automatically manage dependencies. To download them explicitly:

```bash
go mod download
```

### 3. Build the Executable

```bash
go build -o gok-proxy-proxy ./cmd/proxy/
```

This command compiles the proxy server and creates an executable named `gok-proxy-proxy` in the project's root directory.

## Configuration

`gok-proxy` is configured using a `config.yaml` file located in the project root. If the file doesn't exist, it will run with default settings upon the first execution where config is read (though it's best to create it). Refer to `pkg/config/config.go` for all available options and their defaults.

A sample `config.yaml`:

```yaml
ServerAddress: ":8080"
MaxConnections: 10000 # Max connections per IP for the fasthttp server
LogLevel: "info" # Logging level: "debug", "info", "warn", "error"
MaxRequestsPerConn: 5000 # Max requests per connection (server-side)

# Client-side connection pool settings (for outgoing proxied requests)
ClientReadTimeoutSeconds: 15
ClientWriteTimeoutSeconds: 15
ClientMaxIdleConnDurationSeconds: 60
```

## Running gok-proxy

Execute the compiled binary:

```bash
./gok-proxy-proxy
```

The proxy server will start, and log output (by default, to stdout) will indicate its status.

## Testing the Proxy

Use a tool like `curl` to test the proxy server. Replace `127.0.0.1:8080` if your `ServerAddress` is different.

### HTTP Request

```bash
curl -x http://127.0.0.1:8080 http://ifconfig.io
```

### HTTPS Request (via CONNECT tunnel)

```bash
curl -p -x http://127.0.0.1:8080 https://ifconfig.io
```

_(The `-p` flag tells curl to use the CONNECT method for HTTPS proxying.)_

## Project Structure

```plaintext
gok-proxy/
├── cmd/
│   └── proxy/            # Main application package
│       └── main.go       # Application entry point, server setup, and lifecycle management
├── pkg/
│   ├── config/           # Configuration loading and validation
│   │   └── config.go
│   ├── handler/          # HTTP/HTTPS request handling logic
│   │   └── handler.go
│   ├── log/              # Logging setup (using slog)
│   │   └── log.go
│   ├── metrics/          # Prometheus metrics setup
│   │   └── metrics.go
│   ├── pool/             # fasthttp.Client connection pool for outgoing requests
│   │   └── pool.go
│   └── proxy/            # Proxy server core (fasthttp server wrapper)
│       └── proxy.go
├── config.yaml           # Default/sample configuration file
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
└── README.md             # This file
```

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/josephgoksu/gok-proxy/issues) (if your project is hosted and has an issues page).

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

Adherence to Go best practices and the inclusion of relevant tests (if applicable) are appreciated.

## Acknowledgements

- [fasthttp](https://github.com/valyala/fasthttp): For its foundational high-performance HTTP client/server capabilities.
- `slog` (Go standard library): For robust structured logging.
- [viper](https://github.com/spf13/viper): For versatile configuration management.
- [Prometheus](https://prometheus.io/): For industry-standard metrics collection.

## License

Distributed under the MIT License. See `LICENSE` file for more information .
