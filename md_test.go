package main

import (
	"reflect"
	"testing"

	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

func Test_setupMarkdown(t *testing.T) {
	type args struct {
		css   string
		title string
	}
	tests := []struct {
		name string
		args args
		want *bf.HTMLRenderer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setupMarkdown(tt.args.css, tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setupMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_render(t *testing.T) {
	type args struct {
		data  []byte
		css   string
		title string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := render(tt.args.data, tt.args.css, tt.args.title); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("render() = %v, want %v", got, tt.want)
			}
		})
	}
}
