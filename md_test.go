package main

import (
	"reflect"
	"testing"

	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

var markdownTests = []struct {
	name  string
	css   string
	title string
	want  *bf.HTMLRenderer
}{
	{
		name:  "one",
		css:   "assets/wiki.css",
		title: "test page 1",
		want:  bf.NewHTMLRenderer(bf.HTMLRendererParameters{}),
	},
	{
		name:  "two",
		css:   "https://tilde.institute/tilde.css",
		title: "test page 2",
		want:  bf.NewHTMLRenderer(bf.HTMLRendererParameters{}),
	},
}

func Test_setupMarkdown(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got interface{}
			got = setupMarkdown(tt.css, tt.title)
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

// wrapper function to generate the parameters above and
// pass them to the blackfriday library's parsing function
//func Test_render(t *testing.T) {
//return bf.Run(data, bf.WithRenderer(setupMarkdown(css, title)))
//}
