package spa

import (
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// IndexConfig defines the config for the middleware which determines the path to load
// the SPA index file
type IndexConfig struct {
	// Skipper defines a function to skip middleware
	echomiddleware.Skipper

	// This is required to support redirects and branch builds
	DomainName string

	// This can enabled to support serving of branch builds from the s3 bucket
	SubDomainMode bool
}

// IndexWithConfig configure the index middleware
func IndexWithConfig(config IndexConfig) echo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = echomiddleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if config.Skipper(c) {
				return next(c)
			}

			var pathPrefix string

			p := c.Request().URL.Path
			u := c.Request().URL

			if config.SubDomainMode {
				// if we use host it can include the port, hostname:port, whereas hostname is just the hostname
				pathPrefix = extractPathPrefix(config.DomainName, c.Request().URL.Hostname())
			}

			// does the path end in slash?
			if strings.HasSuffix(p, "/") {
				u.Path = path.Join("/", pathPrefix, "index.html")
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
