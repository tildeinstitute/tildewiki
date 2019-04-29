package main

import (
	"net/http"
	"reflect"
	"testing"
)

func Test_pageHandler(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		r        *http.Request
		filename string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageHandler(tt.args.w, tt.args.r, tt.args.filename)
		})
	}
}

func Test_indexHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_iconHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iconHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_cssHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cssHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_validatePath(t *testing.T) {
	type args struct {
		fn func(http.ResponseWriter, *http.Request, string)
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validatePath(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
