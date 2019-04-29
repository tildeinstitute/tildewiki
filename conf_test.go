package main

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func Test_initConfigParams(t *testing.T) {
	tests := []struct {
		name string
		want *regexp.Regexp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initConfigParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initConfigParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_error500(t *testing.T) {
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
			error500(tt.args.w, tt.args.r)
		})
	}
}

func Test_error404(t *testing.T) {
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
			error404(tt.args.w, tt.args.r)
		})
	}
}
