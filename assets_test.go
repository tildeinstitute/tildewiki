package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func Test_setUpUsTheWiki(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUpUsTheWiki()
		})
	}
}

func Test_iconType(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iconType(tt.args.filename); got != tt.want {
				t.Errorf("iconType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cssLocal(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cssLocal(); got != tt.want {
				t.Errorf("cssLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
