package pool

import (
	"sync"
	"time"

	"github.com/josephgoksu/gok-proxy/pkg/config"

	"github.com/valyala/fasthttp"
)

var (
	connPool     sync.Pool
	clientConfig *config.Config
)

// InitConnPool initializes the connection pool with specific client configurations.
// This function should be called once at application startup.
func InitConnPool(cfg *config.Config) {
	clientConfig = cfg
	connPool = sync.Pool{
		New: func() interface{} {
			return &fasthttp.Client{
				ReadTimeout:         time.Duration(clientConfig.ClientReadTimeoutSeconds) * time.Second,
				WriteTimeout:        time.Duration(clientConfig.ClientWriteTimeoutSeconds) * time.Second,
				MaxIdleConnDuration: time.Duration(clientConfig.ClientMaxIdleConnDurationSeconds) * time.Second,
				// Allow HTTP and HTTPS connections by default through fasthttp.Dial
				Dial: fasthttp.Dial,
				// Other potentially useful settings:
				// MaxConnsPerHost: clientConfig.ClientMaxConnsPerHost,
				// MaxResponseBodySize: clientConfig.ClientMaxResponseBodySize,
			}
		},
	}
}

func GetConnection() *fasthttp.Client {
	return connPool.Get().(*fasthttp.Client)
}

func PutConnection(client *fasthttp.Client) {
	connPool.Put(client)
}
