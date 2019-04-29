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

var cssLocalTests = []struct {
	name []byte
	want bool
}{
	{
		name: []byte("https://google.com/test.css"),
		want: false,
	},
	{
		name: []byte("style.css"),
		want: true,
	},
}

func Test_cssLocal(t *testing.T) {
	for _, tt := range cssLocalTests {
		t.Run(string(tt.name), func(t *testing.T) {
			if got := cssLocal(tt.name); got != tt.want {
				t.Errorf("cssLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Benchmark_cssLocal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range cssLocalTests {
			cssLocal(c.name)
		}
	}
}
