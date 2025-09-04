package config

import (
	"net/http"
	"os"
	"time"

	utils "github.com/noo8xl/anvil-gateway/utils/server"
)

type Config struct {
	Services    map[string]string
	TLS         TLSConfig
	RetryConfig RetryConfig
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

type RetryConfig struct {
	MaxAttempts int
	Backoff     time.Duration
}

func GetServerConfig(address string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:           address,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler,
	}
}

// GetaServerAddressByServiceName -> get a dev or prod service address depends on opt value (0 or 1)
func GetServerAddressByServiceName(serviceName string) string {

	opt := os.Getenv("GO_ENV")

	switch opt {
	case "development":
		return utils.GetLocalServerAddress(serviceName)
	case "test":
		return utils.GetLocalServerAddress(serviceName)
	case "production":
		return utils.GetProductionServerAddress(serviceName)
	default:
		return "Got a wrong ENV option name: " + opt
	}
}

func LoadConfig() *Config {
	conf := &Config{
		Services: map[string]string{
			"auth":          GetServerAddressByServiceName("auth"),
			"profile":       GetServerAddressByServiceName("profile"),
			"orders":        GetServerAddressByServiceName("orders"),
			"blog":          GetServerAddressByServiceName("blog"),
			"offers":        GetServerAddressByServiceName("offers"),
			"reviews":       GetServerAddressByServiceName("reviews"),
			"payments":      GetServerAddressByServiceName("payments"),
			"notifications": GetServerAddressByServiceName("notifications"),
		},
		RetryConfig: RetryConfig{
			MaxAttempts: 3,
			Backoff:     1 * time.Second,
		},
	}

	env := os.Getenv("GO_ENV")

	if env == "production" {
		conf.TLS = TLSConfig{
			CertFile: "certs/prod.crt",
			KeyFile:  "certs/prod.key",
		}
	} else {
		// set TLS config for development or testing environment
		conf.TLS = TLSConfig{
			CertFile: "certs/localhost.crt",
			KeyFile:  "certs/localhost.key",
		}
	}

	return conf
}
