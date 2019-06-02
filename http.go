package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
)

// Attach requester's IP address to context value
func newCtxUserIP(ctx context.Context, r *http.Request) context.Context {

	base := strings.Split(r.RemoteAddr, ":")
	uip := base[0]

	if _, ok := r.Header["X-Forwarded-For"]; ok {
		proxied := r.Header["X-Forwarded-For"]
		base = strings.Split(proxied[len(proxied)-1], ":")
		uip = base[0]
	}

	return context.WithValue(ctx, ctxKey, uip)
}

// Retrieve an IP address from context passed with the request
func getIPfromCtx(ctx context.Context) net.IP {

	uip, ok := ctx.Value(ctxKey).(string)
	if !ok {
		log.Printf("Error retrieving IP from request.\n")
	}

	return net.ParseIP(uip)
}

func ipMiddleware(hop http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := newCtxUserIP(r.Context(), r)
		hop.ServeHTTP(w, r.WithContext(ctx))
	})
}
