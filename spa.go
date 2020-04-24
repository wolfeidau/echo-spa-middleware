package spa

import (
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// DefaultIndexFilename the filename used as the default index
const DefaultIndexFilename = "index.html"

// IndexConfig defines the config for the middleware which determines the path to load
// the SPA index file
type IndexConfig struct {
	// Skipper defines a function to skip middleware
	echomiddleware.Skipper

	// This is required to support redirects and branch builds
	DomainName string

	// This can enabled to support serving of branch builds from a folder in a static files store or route
	SubDomainMode bool

	// The name of the file used as the index, defaults to index.html
	IndexFilename string
}

// IndexWithConfig configure the index middleware
func IndexWithConfig(cfg IndexConfig) echo.MiddlewareFunc {

	if cfg.Skipper == nil {
		cfg.Skipper = echomiddleware.DefaultSkipper
	}

	if cfg.IndexFilename == "" {
		cfg.IndexFilename = DefaultIndexFilename
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if cfg.Skipper(c) {
				return next(c)
			}

			var pathPrefix string

			p := c.Request().URL.Path
			u := c.Request().URL

			if cfg.SubDomainMode {
				// if we use host it can include the port, hostname:port, whereas hostname is just the hostname
				pathPrefix = extractPathPrefix(cfg.DomainName, c.Request().URL.Hostname())
			}

			// does the path end in slash?
			if strings.HasSuffix(p, "/") {
				u.Path = path.Join("/", pathPrefix, cfg.IndexFilename)
			} else {
				u.Path = path.Join("/", pathPrefix, p)
			}

			return next(c)
		}
	}
}

func extractPathPrefix(domainName, host string) string {
	if host == domainName {
		return ""
	}

	subdomain := strings.TrimSuffix(host, domainName)

	if !strings.HasSuffix(subdomain, ".") {
		return ""
	}

	return strings.TrimSuffix(subdomain, ".")
}
