package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"
	"time"
)

var loadPageCases = []struct {
	name     string
	filename string
	want     *Page
	wantErr  bool
}{{
	name:     "example.md",
	filename: "pages/example.md",
	want:     &Page{},
	wantErr:  false,
},
	{
		name:     "fake.md",
		filename: "pages/fake.md",
		want:     &Page{},
		wantErr:  true,
	},
}

func Test_loadPage(t *testing.T) {
	for _, tt := range loadPageCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadPage(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("loadPage() = %v, want %v", got, tt.want)
			//}
		})
	}
}
func Benchmark_loadPage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range loadPageCases {
			_, err := loadPage(c.filename)
			if err != nil {
				continue
			}
		}
	}
}

var metaTestBytes, _ = ioutil.ReadFile("pages/example.md")
var metaTestReader = bytes.NewReader(metaTestBytes)
var metaTestScanner = bufio.NewScanner(metaTestReader)
var getMetaCases = []struct {
	name      string
	data      *bufio.Scanner
	titlewant string
	descwant  string
	authwant  string
}{
	{
		name:      "example",
		data:      metaTestScanner,
		titlewant: "Example Page",
		descwant:  "Example page for the wiki",
		authwant:  "`by gbmor`",
	},
}

func Test_getTitle(t *testing.T) {
	for _, tt := range getMetaCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTitle(tt.data); got != tt.titlewant {
				t.Errorf("getTitle() = %v, want %v", got, tt.titlewant)
			}
		})
		metaTestReader.Reset(metaTestBytes)
	}
}
func Benchmark_getTitle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range getMetaCases {
			getTitle(tt.data)
			metaTestReader.Reset(metaTestBytes)
		}
	}
}

func Test_getDesc(t *testing.T) {
	for _, tt := range getMetaCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDesc(tt.data); got != tt.descwant {
				t.Errorf("getDesc() = %v, want %v", got, tt.descwant)
			}
		})
		metaTestReader.Reset(metaTestBytes)
	}
}
func Benchmark_getDesc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range getMetaCases {
			getDesc(tt.data)
			metaTestReader.Reset(metaTestBytes)
		}
	}
}

func Test_getAuthor(t *testing.T) {
	for _, tt := range getMetaCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAuthor(tt.data); got != tt.authwant {
				t.Errorf("getAuthor() = %v, want %v", got, tt.authwant)
			}
		})
		metaTestReader.Reset(metaTestBytes)
	}
}
func Benchmark_getAuthor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range getMetaCases {
			getAuthor(tt.data)
			metaTestReader.Reset(metaTestBytes)
		}
	}
}

var genIndexCases = []struct {
	name string
}{
	{
		name: "index",
	},
}

func Test_genIndex(t *testing.T) {
	for _, tt := range genIndexCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := genIndex(); len(got) <= 0 {
				t.Errorf("genIndex(), got %v bytes.", got)
			}
		})
	}
}
func Benchmark_genIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range genIndexCases {
			out := genIndex()
			if len(out) <= 0 {
				continue
			}
		}
	}
}

var tallyPagesCases = []struct {
	name string
}{
	{
		name: "index",
	},
}

func Test_tallyPages(t *testing.T) {
	for _, tt := range tallyPagesCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tallyPages(); len(got) <= 0 {
				t.Errorf("tallyPages() = %v", got)
			}
		})
	}
}
func Benchmark_tallyPages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range tallyPagesCases {
			out := tallyPages()
			if len(out) <= 0 {
				continue
			}
		}
	}
}

type fields struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       []byte
}

var PageCacheCases = []struct {
	name    string
	fields  fields
	wantErr bool
}{
	{
		name: "test1",
		fields: fields{
			Longname:  "pages/test1.md",
			Shortname: "test1.md",
		},
		wantErr: false,
	},
	{
		name: "example",
		fields: fields{
			Longname:  "pages/example.md",
			Shortname: "example.md",
		},
		wantErr: false,
	},
}

func TestPage_cache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
				Title:     tt.fields.Title,
				Desc:      tt.fields.Desc,
				Author:    tt.fields.Author,
				Modtime:   tt.fields.Modtime,
				Body:      tt.fields.Body,
				Raw:       tt.fields.Raw,
			}
			if err := page.cache(); (err != nil) != tt.wantErr {
				t.Errorf("Page.cache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Benchmark_Page_cache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range PageCacheCases {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			err := page.cache()
			if err != nil {
				continue
			}
		}
	}
}

func TestPage_checkCache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
				Title:     tt.fields.Title,
				Desc:      tt.fields.Desc,
				Author:    tt.fields.Author,
				Modtime:   tt.fields.Modtime,
				Body:      tt.fields.Body,
				Raw:       tt.fields.Raw,
			}
			err := page.cache()
			if err != nil {
				t.Errorf("Page.checkCache() returned error: %v\n", err)
			}
			var got interface{} = page.checkCache()
			switch got.(type) {
			case bool:
				return
			default:
				t.Errorf("Page.checkCache() = %v", got)
			}
		})
	}
}
func Benchmark_Page_checkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range PageCacheCases {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			err := page.cache()
			if err != nil {
				continue
			}
			var maybe interface{} = page.checkCache()
			if maybe.(bool) {
				continue
			}
		}
	}
}

func Test_genPageCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "first",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genPageCache()
		})
	}
}
func Benchmark_genPageCache(b *testing.B) {
	tests := []struct {
		name string
	}{
		{
			name: "first",
		},
	}
	for i := 0; i < b.N; i++ {
		for range tests {
			genPageCache()
		}
	}
}
