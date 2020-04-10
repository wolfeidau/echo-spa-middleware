# echo-spa-middleware

This middleware is specifically designed to resolve the `index.html` given a few parameters and is typical of the logic used to locate and serve an SPA index files.

# Configuration

This echo middleware has a couple of key configuration options.

* **DomainName** - This option enabled the middleware to recognise sub domain request and is used to support branch builds.
* **SubDomainMode** - This can be enabled to support rewriting of sub domain to path prefix.

# Sub Domains

This middleware supports using the sub domain which prefixes `domainName` and use it to form a part of the logic to locate the `indexFile` to be served.

So with a configuration of the following:

```go
e := echo.New()
e.Pre(echomiddleware.AddTrailingSlash()) // required to ensure trailing slash is appended
e.Use(spa.HTTPSRedirectWithConfig(spa.RedirectConfig{
  DomainName: "www.example.com",
  SubDomainMode: true,
}))
```

## branch mode

This middleware supports enabling a branch mode, to support branch builds of an SPA.

At least two DNS records pointing at this service being:

* `www.example.com`
* `*.www.example.com`

Request URI | path served
--- | ---
http://example.com | /index.html
http://mybranch_build.example.com | /mybranch_build/index.html
http://example.com/someroute  | /index.html
http://mybranch_build.example.com/someroute | /mybranch_build/index.html
http://mybranch_build.example.com/img/logo.png | /mybranch_build/img/logo.png

# License

This library is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).