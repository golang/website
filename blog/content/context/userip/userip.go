// Package userip provides functions for extracting a user IP address from a
// request and associating it with a Context.
package userip

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"code.google.com/p/go.net/context"
)

// FromRequest extracts the user IP address from req, if present.
func FromRequest(req *http.Request) (net.IP, error) {
	s := strings.SplitN(req.RemoteAddr, ":", 2)
	userIP := net.ParseIP(s[0])
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// The address &key is the context key.
var key struct{}

// NewContext returns a new Context carrying userIP.
func NewContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, &key, userIP)
}

// FromContext extracts the user IP address from ctx, if present.
func FromContext(ctx context.Context) (net.IP, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	userIP, ok := ctx.Value(&key).(net.IP)
	return userIP, ok
}
