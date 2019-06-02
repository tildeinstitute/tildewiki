package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

// This is a pretty strict test. Make sure the
// output of pageHandler is byte-for-byte what
// I'm expecting it to be.
/*
func Test_pageHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "example",
		},
		{
			name: "test1",
		},
	}

	hush, _ := os.Open("/dev/null")
	log.SetOutput(hush)
	initConfigParams()
	genPageCache()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "localhost:8080/w/"+tt.name, nil)
			pageHandler(w, req)
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				t.Errorf("pageHandler(): %v\n", resp.StatusCode)
			}
			if !bytes.Equal(body, cachedPages[tt.name+".md"].Body) {
				t.Errorf("pageHandler(): Byte mismatch\n")
			}
		})
	}
}
*/

// This is the same test type as pageHandler
func Test_indexHandler(t *testing.T) {
	name := "Index Handler Test"
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080", nil)
	t.Run(name, func(t *testing.T) {
		indexHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("indexHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, indexCache.Body) {
			t.Errorf("indexHandler(): Byte mismatch\n")
		}
	})
}

// This is the same test type as pageHandler
func Test_iconHandler(t *testing.T) {
	name := "Icon Handler Test"
	initConfigParams()

	confVars.mu.RLock()
	icon, _ := ioutil.ReadFile(confVars.assetsDir + "/" + confVars.iconPath)
	confVars.mu.RUnlock()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/icon", nil)
	t.Run(name, func(t *testing.T) {
		iconHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("iconHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, icon) {
			t.Errorf("iconHandler(): Byte mismatch\n")
		}
	})
}

// This is the same test type as pageHandler
func Test_cssHandler(t *testing.T) {
	name := "CSS Handler Test"
	initConfigParams()
	if !cssLocal([]byte(confVars.cssPath)) {
		t.Skipf("cssHandler(): Set to use remote CSS in config, skipping test ...\n")
	}
	css, _ := ioutil.ReadFile(confVars.cssPath)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/css", nil)
	t.Run(name, func(t *testing.T) {
		cssHandler(w, r)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			t.Errorf("cssHandler(): %v\n", resp.StatusCode)
		}
		if !bytes.Equal(body, css) {
			t.Errorf("cssHandler(): Byte mismatch\n")
		}
	})
}

// Tests if /500 returns a status 200, which means
// the handler is working. Doesn't test for 500-triggering
// situations yet.
func Test_error500(t *testing.T) {
	name := "Error 500 Handler Test"
	initConfigParams()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/500", nil)
	t.Run(name, func(t *testing.T) {
		error500(w, r)
		resp := w.Result()
		if resp.StatusCode != 200 {
			t.Errorf("error500(): %v\n", resp.StatusCode)
		}
	})
}

// Tests for a 200 status code because it serves requests
// that fail the regex path validation, rather than a traditional
// 404 status code.
func Test_error404(t *testing.T) {
	name := "Error 404 Handler Test"
	initConfigParams()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080"+confVars.viewPath+"?@$#$", nil)
	t.Run(name, func(t *testing.T) {
		error404(w, r)
		resp := w.Result()
		if resp.StatusCode != 200 {
			t.Errorf("error404(): %v\n", resp.StatusCode)
		}
	})
}
