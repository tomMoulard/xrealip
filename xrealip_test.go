package xrealip

import (
	"context"
	"net/http"
	"testing"
)

// TODO test ranges

const (
	xForwardedFor = "X-Forwarded-For"
	xRealIP       = "X-Real-Ip"
)

func TestConfig_getXRealIp(t *testing.T) {
	remoteAddr := "42.42.42.42"
	tests := []struct {
		name    string
		config  Config
		headers http.Header
		want    string
	}{
		{
			name: "empty",
			want: remoteAddr,
		},
		{
			name: "no from",
			config: Config{
				header: xForwardedFor,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "no matchin header",
			config: Config{
				header: xForwardedFor,
			},
			want: remoteAddr,
		},
		{
			name: "from but no match",
			config: Config{
				from:   []string{"127.0.0.2"},
				header: xForwardedFor,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "from and match",
			config: Config{
				from:   []string{remoteAddr},
				header: xForwardedFor,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "from and match multiple address in header",
			config: Config{
				from:   []string{remoteAddr},
				header: xForwardedFor,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1, 127.0.0.2"},
			},
			want: "127.0.0.2",
		},
		{
			name: "from and match multiple address in header 2",
			config: Config{
				from:   []string{remoteAddr},
				header: xForwardedFor,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.2, 127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "no header value",
			config: Config{
				from:   []string{remoteAddr},
				header: xForwardedFor,
			},
			want: remoteAddr,
		},
		{
			name: "XRIP set",
			config: Config{
				from:   []string{remoteAddr},
				header: xRealIP,
			},
			headers: http.Header{
				xRealIP: []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "XRIP set with from",
			config: Config{
				from:   []string{"10.0.0.0"},
				header: xRealIP,
			},
			headers: http.Header{
				xRealIP: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "recursive: no from",
			config: Config{
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "recursive: no matchin header",
			config: Config{
				header:    xForwardedFor,
				recursive: true,
			},
			want: remoteAddr,
		},
		{
			name: "recursive: from but no match",
			config: Config{
				from:      []string{"127.0.0.2"},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "recursive: from and match",
			config: Config{
				from:      []string{remoteAddr},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "recursive: from and match multiple address in header",
			config: Config{
				from:      []string{remoteAddr},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.1, 127.0.0.2"},
			},
			want: "127.0.0.2",
		},
		{
			name: "recursive: from and match multiple address in header 2",
			config: Config{
				from:      []string{remoteAddr},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.2, 127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "recursive: from 2 addr and match multiple address in header 2",
			config: Config{
				from:      []string{remoteAddr, "127.0.0.1"},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.2, 127.0.0.1"},
			},
			want: "127.0.0.2",
		},
		{
			name: "recursive: from 2 addr and match multiple address in header 3",
			config: Config{
				from:      []string{remoteAddr, "127.0.0.1"},
				header:    xForwardedFor,
				recursive: true,
			},
			headers: http.Header{
				xForwardedFor: []string{"127.0.0.3, 127.0.0.2, 127.0.0.1"},
			},
			want: "127.0.0.2",
		},
		{
			name: "recursive: no header value",
			config: Config{
				from:      []string{remoteAddr},
				header:    xForwardedFor,
				recursive: true,
			},
			want: remoteAddr,
		},
		{
			name: "recursive: XRIP set",
			config: Config{
				from:      []string{remoteAddr},
				header:    xRealIP,
				recursive: true,
			},
			headers: http.Header{
				xRealIP: []string{"127.0.0.1"},
			},
			want: "127.0.0.1",
		},
		{
			name: "recursive: XRIP set with from",
			config: Config{
				from:      []string{"10.0.0.0"},
				header:    xRealIP,
				recursive: true,
			},
			headers: http.Header{
				xRealIP: []string{"127.0.0.1"},
			},
			want: remoteAddr,
		},
		{
			name: "recursive: values in header already known",
			config: Config{
				from:      []string{remoteAddr},
				header:    xRealIP,
				recursive: true,
			},
			headers: http.Header{
				xRealIP: []string{remoteAddr},
			},
			want: remoteAddr,
		},
		{
			name: "recursive: multiple values in header already known",
			config: Config{
				from:      []string{remoteAddr},
				header:    xRealIP,
				recursive: true,
			},
			headers: http.Header{
				xRealIP: []string{remoteAddr, "127.0.0.1"},
			},
			want: remoteAddr,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if len(test.headers) == 0 {
				test.headers = make(map[string][]string)
			}
			req := http.Request{Header: test.headers, RemoteAddr: remoteAddr}
			handler := http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
				if got := req.Header.Get(xRealIP); got != test.want {
					t.Errorf("getXRealIP() = %v, want %v", got, test.want)
				}
			})
			d, err := New(context.Background(), handler, &test.config, "test-plugin")
			if err != nil {
				t.Fail()
			}

			d.ServeHTTP(nil, &req)
		})
	}
}
