package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

var setUpUsTheWikiTests = []struct {
	name string
}{
	{
		name: "this test isn't needed but I'm doing it for completeness",
	},
}

func Test_setUpUsTheWiki(t *testing.T) {
	for _, tt := range setUpUsTheWikiTests {
		t.Run(tt.name, func(t *testing.T) {
			setUpUsTheWiki()
		})
	}
}
func Benchmark_setUpUsTheWiki(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range setUpUsTheWikiTests {
			setUpUsTheWiki()
		}
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
