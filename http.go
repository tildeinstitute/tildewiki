package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

func log200(r *http.Request) {
	useragent := r.Header["User-Agent"]
	uip := getIPfromCtx(r.Context())
	log.Printf("**** %v :: 200 :: %v %v :: %v\n", uip, r.Method, r.URL, useragent)
}

// wrapper for testing 500 pages via /500
func error500(w http.ResponseWriter, r *http.Request) {
	log500(w, r, fmt.Errorf("500 Page Accessed Directly, No Error"))
}

// this is a custom 500 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func log500(w http.ResponseWriter, r *http.Request, topErr error) {
	useragent := r.Header["User-Agent"]
	uip := getIPfromCtx(r.Context())
	log.Printf("**** %v :: 500 :: %v %v :: %v :: %v\n", uip, r.Method, r.URL, useragent, topErr.Error())

	confVars.mu.RLock()
	e500 := confVars.assetsDir + "/500.md"
	confVars.mu.RUnlock()

	file, err := ioutil.ReadFile(e500)
	if err != nil {
		log.Printf("Tried to read 500.md: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "500: Internal Server Error"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is a custom 404 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error404(w http.ResponseWriter, r *http.Request) {
	confVars.mu.RLock()
	e404 := confVars.assetsDir + "/404.md"
	confVars.mu.RUnlock()

	file, err := ioutil.ReadFile(e404)
	if err != nil {
		log.Printf("Tried to read 404.md: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "404: Not Found"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err.Error())
		error500(w, r)
	}
}
