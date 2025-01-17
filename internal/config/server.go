package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sumup-oss/go-pkgs/errors"
)

type ServerConfig struct {
	// GracefulShutdownTimeout is a deadline for all services to shut down upon initiated graceful shutdown.
	// The graceful shutdown is initiated with SIGINT or SIGTERM.
	GracefulShutdownTimeout time.Duration `default:"60s" envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`

	// ListenAddress is the IP address and port where the metrics HTTP server listens at.
	ListenAddress string `default:"0.0.0.0:8080" envconfig:"LISTEN_ADDRESS"`

	// ReadTimeout is the maximum duration for reading the entire request, including the body.
	ReadTimeout time.Duration `default:"15s" envconfig:"READ_TIMEOUT"`

	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	ReadHeaderTimeout time.Duration `default:"15s" envconfig:"READ_HEADER_TIMEOUT"`

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration `default:"15s" envconfig:"WRITE_TIMEOUT"`

	// MaxHeaderBytes is the maximum number of bytes the server will read parsing the request header's keys and values.
	MaxHeaderBytes int `default:"1000000" envconfig:"MAX_HEADER_BYTES"`

	// CorsAllowedOrigins is a comma-separated list of origins allowed via CORS.
	CorsAllowedMethods []string `default:"GET,POST,DELETE,OPTIONS" envconfig:"CORS_ALLOWED_METHODS"`

	// CorsAllowedOrigins is a comma-separated list of origins allowed via CORS.
	CorsAllowedOrigins []string `default:"http://*,https://*" envconfig:"CORS_ALLOWED_ORIGINS"`

	// CorsAllowedHeaders is a comma-separated list of headers allowed via CORS.
	CorsAllowedHeaders []string `default:"X-PINGOTHER,Accept,Authorization,Content-Type,X-CSRF-Token" envconfig:"CORS_ALLOWED_HEADERS"` //nolint:lll

	// CorsExposedHeaders is a comma-separated list of headers exposed via CORS.
	CorsExposedHeaders []string `default:"Link" envconfig:"CORS_EXPOSED_HEADERS"`

	// CorsAllowCredentials is a boolean value indicating whether the resource allows credentials.
	CorsAllowCredentials bool `default:"true" envconfig:"CORS_ALLOW_CREDENTIALS"`

	// CorsMaxAge is a positive number of seconds indicating how long the results of a preflight request can be cached.
	CorsMaxAge int `default:"300" envconfig:"CORS_MAX_AGE"`

	// CorsDebug enables debug mode for CORS.
	CorsDebug bool `default:"false" envconfig:"CORS_DEBUG"`

	Log       Log
	Database  Database
	JWTConfig JWTConfig
}

func NewServerConfig() (*ServerConfig, error) {
	var result ServerConfig

	err := envconfig.Process("", &result)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"failed to process the config environment variables",
		)
	}

	return &result, nil
}
