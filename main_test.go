package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func serve(handler fasthttp.RequestHandler, req *http.Request) (*http.Response, error) {
	ln := fasthttputil.NewInmemoryListener()
	defer ln.Close()

	go func() {
		err := fasthttp.Serve(ln, handler)
		if err != nil {
			panic(fmt.Errorf("failed to serve: %v", err))
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return ln.Dial()
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return client.Do(req)
}

func doRequest(t *testing.T, env string, srcURI string, redURI string) {

	os.Setenv("DOMAIN_1", env)
	os.Setenv("HEALTH", "/url-redirect/health")
	loadEnv()

	r, err := http.NewRequest("GET", srcURI, nil)
	if err != nil {
		t.Error(err)
	}

	res, err := serve(fastHTTPHandler, r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res.StatusCode, 302)
	assert.Equal(t, res.Header.Get("Location"), redURI)
}

// TestDomainRedirection method
func TestDomainRedirection(t *testing.T) {
	doRequest(t, "test;redirect-domain;;", "http://test", "http://redirect-domain/")
}

// TestDomainRedirectionWithPathPrefix method
func TestDomainRedirectionWithPathPrefix(t *testing.T) {
	doRequest(t, "test;redirect-domain;add_path;", "http://test", "http://redirect-domain/add_path/")
}

// TestDomainRedirectionWithParameter method
func TestDomainRedirectionWithParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;;add=parameter", "http://test", "http://redirect-domain/?add=parameter")
}

// TestDomainRedirectionWithPathAndParameter method
func TestDomainRedirectionWithPathAndParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;add-path;add=parameter", "http://test", "http://redirect-domain/add-path/?add=parameter")
}

// TestDomainRedirectionWithPathAndExistParameter method
func TestDomainRedirectionWithPathAndExistParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;add-path;add=parameter", "http://test/?exist=param", "http://redirect-domain/add-path/?exist=param&add=parameter")
}

// TestDomainRedirectionDocument method
func TestDomainRedirectionDocument(t *testing.T) {
	doRequest(t, "test;redirect-domain;;", "http://test/docs", "http://redirect-domain/docs")
}

// TestDomainRedirectionDocumentAddPath method
func TestDomainRedirectionDocumentAddPath(t *testing.T) {
	doRequest(t, "test;redirect-domain;add-path;", "http://test/docs", "http://redirect-domain/add-path/docs")
}

// TestDomainRedirectionDocumentAddPathParameter method
func TestDomainRedirectionDocumentAddPathParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;add-path;add=parameter", "http://test/docs", "http://redirect-domain/add-path/docs?add=parameter")
}

// TestDomainRedirectionDocumentAddParameter method
func TestDomainRedirectionDocumentAddParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;;add=parameter", "http://test/docs", "http://redirect-domain/docs?add=parameter")
}

// TestDomainRedirectionDocumentAddExistingParameter method
func TestDomainRedirectionDocumentAddExistingParameter(t *testing.T) {
	doRequest(t, "test;redirect-domain;;add=parameter", "http://test/docs?exist=param", "http://redirect-domain/docs?exist=param&add=parameter")
}

func TestHealthCheck(t *testing.T) {

	r, err := http.NewRequest("GET", "http://test/url-redirect/health", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := serve(fastHTTPHandler, r)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res.StatusCode, 200)
}
