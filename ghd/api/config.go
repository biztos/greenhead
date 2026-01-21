package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config defines the configuration of the API.
//
// It is included as part of the Runner config.
//
// Fields beginning with "App" are passed to the Fiber app.
type Config struct {

	// General configuration:
	ListenAddress string `toml:"listen_address"` // Address for serving, e.g. ":3000"
	LogFiber      bool   `toml:"log_fiber"`      // Use default Fiber logger for requests.

	// Access control:
	Roles      []*Role `toml:"roles"`       // Roles defining access.
	Keys       []*Key  `toml:"keys"`        // Keys mapping to roles by name.
	AccessFile string  `toml:"access_file"` // TOML file for (more) Roles and Keys.
	RawKeys    bool    `toml:"raw_keys"`    // Use raw, unencoded API keys.
	NoKeys     bool    `toml:"no_keys"`     // DO NOT require API keys.
	NoUI       bool    `toml:"no_ui"`       // DO NOT expose the web UI.

	// Fiber app config; see Fiber docs for specifics:
	AppPrefork                 bool          `toml:"app_prefork"`
	AppServerHeader            string        `toml:"app_server_header"`
	AppBodyLimit               int           `toml:"app_body_limit"`
	AppConcurrency             int           `toml:"app_concurrency"`
	AppReadTimeout             time.Duration `toml:"app_read_timeout"`
	AppWriteTimeout            time.Duration `toml:"app_write_timeout"`
	AppIdleTimeout             time.Duration `toml:"app_idle_timeout"`
	AppProxyHeader             string        `toml:"app_proxy_header"`
	AppDisableStartupMessage   bool          `toml:"app_disable_startup_message"`
	AppEnableTrustedProxyCheck bool          `toml:"app_enable_trusted_proxy_check"`
	AppTrustedProxies          []string      `toml:"app_trusted_proxies"`
}

// FiberConfig returns a fiber.Config from the App* fields in cfg, with
// AppName set to ident.
func (cfg *Config) FiberConfig(ident string) fiber.Config {

	// Cf: https://pkg.go.dev/github.com/gofiber/fiber/v2#Config
	return fiber.Config{
		AppName:                 ident,
		Prefork:                 cfg.AppPrefork,
		ServerHeader:            cfg.AppServerHeader,
		BodyLimit:               cfg.AppBodyLimit,
		Concurrency:             cfg.AppConcurrency,
		ReadTimeout:             cfg.AppReadTimeout,
		WriteTimeout:            cfg.AppWriteTimeout,
		IdleTimeout:             cfg.AppIdleTimeout,
		ProxyHeader:             cfg.AppProxyHeader,
		DisableStartupMessage:   cfg.AppDisableStartupMessage,
		EnableTrustedProxyCheck: cfg.AppEnableTrustedProxyCheck,
		TrustedProxies:          cfg.AppTrustedProxies,
		// EnablePrintRoutes ?
	}
}
