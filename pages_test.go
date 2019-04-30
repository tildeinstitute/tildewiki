package main

import (
	"reflect"
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

func Test_getTitle(t *testing.T) {
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
			if got := getTitle(tt.args.filename); got != tt.want {
				t.Errorf("getTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDesc(t *testing.T) {
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
			if got := getDesc(tt.args.filename); got != tt.want {
				t.Errorf("getDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAuthor(t *testing.T) {
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
			if got := getAuthor(tt.args.filename); got != tt.want {
				t.Errorf("getAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genIndex(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genIndex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tallyPages(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tallyPages(); got != tt.want {
				t.Errorf("tallyPages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPage_cache(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
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

func TestPage_checkCache(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
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
			if got := page.checkCache(); got != tt.want {
				t.Errorf("Page.checkCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genPageCache(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			genPageCache()
		})
	}
}
