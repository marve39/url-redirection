package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/valyala/fasthttp"
)

var domainRedirection map[string][]string
var healthCheckPath string

// loadEnv method
// DOMAIN_1=source_domain;redirect_domain;path-prefix;add-parameter)
func loadEnv() {
	domainRedirection = make(map[string][]string)
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		if strings.HasPrefix(variable[0], "DOMAIN_") {
			tempDomain := strings.Split(strings.Join(variable[1:], "="), ";")
			if len(tempDomain) >= 4 {
				domain := strings.ToLower(tempDomain[0])
				domainRedirection[domain] = tempDomain
				log.Info(fmt.Sprintf("Load Domain Config [%s] => %s", domain, strings.Join(tempDomain, ";")))
			}
		}
		if strings.ToLower(variable[0]) == "health" {
			healthCheckPath = variable[1]
		}
	}
}

// fastHTTPHandler method
// DOMAIN_1=test;redirect-domain;path_prefix;append_parameter
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	inboundURL := string(ctx.Request.URI().RequestURI())
	if inboundURL == healthCheckPath {
		ctx.Response.SetBodyString("OK")
		return
	}

	domain := string(ctx.URI().Host())
	originalURL := string(ctx.Request.URI().FullURI())
	redirectURL := originalURL

	if record, ok := domainRedirection[domain]; ok {
		destDomain := record[1]
		addPathPrefix := record[2]
		addParameter := record[3]

		ctx.Request.URI().SetHost(destDomain)
		if len(addPathPrefix) > 0 {
			ctx.Request.URI().SetPath(fmt.Sprintf("%s/%s", addPathPrefix, string(ctx.Request.URI().Path())))
		}

		if len(addParameter) > 0 {
			existingParam := ctx.Request.URI().QueryArgs().String()
			if len(existingParam) > 0 {
				ctx.Request.URI().SetQueryString(fmt.Sprintf("%s&%s", existingParam, addParameter))
			} else {
				ctx.Request.URI().SetQueryString(addParameter)
			}
		}
		redirectURL = string(ctx.Request.URI().FullURI())
		ctx.Redirect(redirectURL, 302)
	}

	log.Info(fmt.Sprintf("%s -> %s", originalURL, redirectURL))
}

func main() {
	loadEnv()
	log.Info("Listening to :80")
	err := fasthttp.ListenAndServe(":80", fastHTTPHandler)
	if err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}
