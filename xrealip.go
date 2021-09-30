// Package xrealip a demo plugin.
package xrealip

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
)

// Config the plugin configuration.
// See https://nginx.org/en/docs/http/ngx_http_realip_module.html
type Config struct {
	// from is the equivalent as set_real_ip_from from nginx
	// It defines trusted addresses that are known to send correct replacement
	// addresses. If the special value unix: is specified, all UNIX-domain
	// sockets will be trusted. Trusted addresses may also be specified using a
	// hostname (1.13.1).
	from []string
	// header is the equivalent as real_ip_header from nginx
	// It defines the request header field whose value will be used to
	// replace the client address.
	header string
	// recursive is the equivalent as real_ip_recursive from nginx
	// If recursive search is disabled, the original client address
	// that matches one of the trusted addresses is replaced by the last address
	// sent in the request header field defined by the real_ip_header directive.
	// If recursive search is enabled, the original client address that matches
	// one of the trusted addresses is replaced by the last non-trusted address
	// sent in the request header field.
	recursive bool
}

// CreateConfig populates the Config data object.
func CreateConfig() *Config {
	return &Config{}
}

// Demo a Demo plugin.
type Demo struct {
	config *Config
	ctx    context.Context
	fromIP []net.IP
	name   string
	next   http.Handler
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Default().Printf("xrealip plugin, loading configuration: %+v", config)

	var fromIP []net.IP
	for _, ip := range config.from {
		fromIP = append(fromIP, net.ParseIP(ip))
	}

	return &Demo{
		config: config,
		ctx:    ctx,
		fromIP: fromIP,
		name:   name,
		next:   next,
	}, nil
}

func (a *Demo) trustRemote(remoteAddr net.IP) bool {
	for _, ip := range a.fromIP {
		if ip.Equal(remoteAddr) {
			return true
		}
	}

	return false
}

func (a *Demo) lastNotMatched(headerValues []string) string {
	for i := len(headerValues) - 1; i >= 0; i-- {
		value := net.ParseIP(strings.TrimSpace(headerValues[i]))

		var matched bool
		for _, ip := range a.fromIP {
			if ip.Equal(value) {
				matched = true

				break
			}
		}
		if !matched {
			return strings.TrimSpace(headerValues[i])
		}
	}
	return strings.TrimSpace(headerValues[len(headerValues)-1])
}

func (a *Demo) getXRealIP(req *http.Request) string {
	headerValue := req.Header.Get(a.config.header)
	if headerValue == "" {
		return req.RemoteAddr
	}

	headerValues := strings.Split(headerValue, ",")
	remoteAddrIP := net.ParseIP(req.RemoteAddr)
	if a.trustRemote(remoteAddrIP) {
		if !a.config.recursive {
			return strings.TrimSpace(headerValues[len(headerValues)-1])
		}

		return a.lastNotMatched(headerValues)
	}

	return req.RemoteAddr
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.Header.Set("X-Real-Ip", a.getXRealIP(req))
	a.next.ServeHTTP(rw, req)
}
