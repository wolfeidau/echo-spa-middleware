package spa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_pathPrefix(t *testing.T) {

	assert := require.New(t)

	type args struct {
		domainName string
		host       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should match and return valid subdomain",
			args: args{
				domainName: "app.example.com",
				host:       "customer.app.example.com",
			},
			want: "customer",
		},
		{
			name: "should skip invalid subdomain and return empty string",
			args: args{
				domainName: "app.example.com",
				host:       "customerapp.example.com",
			},
			want: "",
		},
		{
			name: "should skip unrelated subdomain and return empty string",
			args: args{
				domainName: "app.example.com",
				host:       "customer.app.beer.com",
			},
			want: "",
		},
		{
			name: "should skip empty subdomain and return empty string",
			args: args{
				domainName: "app.example.com",
				host:       ".app.example.com",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(tt.want, extractPathPrefix(tt.args.domainName, tt.args.host))
		})
	}
}

func TestIndexSubdomain(t *testing.T) {

	assert := require.New(t)

	e := echo.New()

	e.Use(IndexWithConfig(IndexConfig{
		DomainName:    "app.example.com",
		SubDomainMode: true,
	}))

	tests := []struct {
		name     string
		url      string
		wantPath string
	}{
		{
			name:     "should rewrite with subdomain prefix",
			url:      "http://sup.app.example.com/add-slash/",
			wantPath: "/sup/index.html",
		},
		{
			name:     "should rewrite with no prefix",
			url:      "http://app.example.com/add-slash/",
			wantPath: "/index.html",
		},
		{
			name:     "should rewrite mismatched domain with no prefix",
			url:      "http://not.example.com/add-slash/",
			wantPath: "/index.html",
		},
		{
			name:     "should not rewrite with no trailing slash",
			url:      "http://app.example.com/add-slash.gif",
			wantPath: "/add-slash.gif",
		},
		{
			name:     "should not rewrite subdomain prefix with no trailing slash",
			url:      "http://customer.app.example.com/add-slash.gif",
			wantPath: "/customer/add-slash.gif",
		},
		{
			name:     "should not rewrite subdomain prefix with no trailing slash",
			url:      "http://customer.app.example.com/add-slash/",
			wantPath: "/customer/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, tt.url, nil)
			assert.NoError(err)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)
			assert.Equal(tt.wantPath, req.URL.Path)
		})
	}

}
