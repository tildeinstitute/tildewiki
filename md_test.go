package main

import (
	"io/ioutil"
	"reflect"
	"testing"

	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

var mdTestData1, _ = ioutil.ReadFile("pages/example.md")
var mdTestData2, _ = ioutil.ReadFile("pages/test1.md")
var markdownTests = []struct {
	name  string
	css   string
	title string
	data  []byte
	want  *bf.HTMLRenderer
}{
	{
		name:  "one",
		css:   "assets/wiki.css",
		title: "test page 1",
		data:  mdTestData1,
		want:  bf.NewHTMLRenderer(bf.HTMLRendererParameters{}),
	},
	{
		name:  "two",
		css:   "https://tilde.institute/tilde.css",
		title: "test page 2",
		data:  mdTestData2,
		want:  bf.NewHTMLRenderer(bf.HTMLRendererParameters{}),
	},
}

func Test_setupMarkdown(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got interface{} = setupMarkdown(tt.css, tt.title)
			if _, ok := got.(*bf.HTMLRenderer); !ok {
				t.Errorf("setupMarkdown() returned incorrect type: %v", reflect.TypeOf(got))
			}
		})
	}
}
func Benchmark_setupMarkdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			setupMarkdown(c.css, c.title)
		}
	}
}

func Test_render(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got interface{} = render(tt.data, tt.css, tt.title)
			if _, ok := got.([]byte); !ok {
				t.Errorf("render() didn't return byte array: %v", reflect.TypeOf(got))
			}
		})
	}
}
func Benchmark_render(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			render(c.data, c.css, c.title)
		}
	}
}
