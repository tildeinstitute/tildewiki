package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// to quiet the function output during
// testing and benchmarks
var hush, _ = os.Open("/dev/null")

var buildPageCases = []struct {
	name     string
	filename string
	want     *Page
	wantErr  bool
}{
	{
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

func Test_buildPage(t *testing.T) {
	log.SetOutput(hush)
	for _, tt := range buildPageCases {
		t.Run(tt.name, func(t *testing.T) {
			testpage, err := buildPage(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildPage() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if testpage == nil && !tt.wantErr {
				t.Errorf("buildPage() returned nil bytes when it wasn't expected.\n")
			}
		})
	}
}
func Benchmark_buildPage(b *testing.B) {
	log.SetOutput(hush)
	for i := 0; i < b.N; i++ {
		for _, c := range buildPageCases {
			_, err := buildPage(c.filename)
			if (err != nil) != c.wantErr {
				b.Errorf("buildPage benchmark failed: %v\n", err)
			}
		}
	}
}

var metaBytes, _ = ioutil.ReadFile("pages/example.md")
var metaTestBytes pagedata = metaBytes
var getMetaCases = []struct {
	name      string
	data      pagedata
	titlewant string
	descwant  string
	authwant  string
}{
	{
		name:      "example",
		data:      metaTestBytes,
		titlewant: "Example Page",
		descwant:  "Example page for the wiki",
		authwant:  "gbmor",
	},
}

func Test_getMeta(t *testing.T) {
	for _, tt := range getMetaCases {
		t.Run(tt.name, func(t *testing.T) {
			if title, desc, auth := tt.data.getMeta(); title != tt.titlewant || desc != tt.descwant || auth != tt.authwant {
				t.Errorf("getMeta() = %v, %v, %v .. want %v, %v, %v", title, desc, auth, tt.titlewant, tt.descwant, tt.authwant)
			}
		})
	}
}
func Benchmark_getMeta(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range getMetaCases {
			tt.data.getMeta()
		}
	}
}

func Test_genIndex(t *testing.T) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	t.Run("genIndex() test", func(t *testing.T) {
		if got := genIndex(); got == nil {
			t.Errorf("genIndex(), got %v bytes.", got)
		}
	})
}
func Benchmark_genIndex(b *testing.B) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	for i := 0; i < b.N; i++ {
		indexCache.Modtime = time.Time{}
		genIndex()
	}
}

var tallyPagesPagelist = make([]byte, 0, 1)
var tallyPagesBuf = bytes.NewBuffer(tallyPagesPagelist)

// Currently tests for whether the buffer is being written to.
// Also checks if the anchor tag was replaced in the buffer.
func Test_tallyPages(t *testing.T) {
	t.Run("tallyPages test", func(t *testing.T) {
		if tallyPages(tallyPagesBuf); tallyPagesBuf == nil {
			t.Errorf("tallyPages() wrote nil to buffer\n")
		}
		bufscan := bufio.NewScanner(tallyPagesBuf)
		for bufscan.Scan() {
			if bufscan.Text() == "<!--pagelist-->" {
				t.Errorf("tallyPages() - Did not replace anchor tag with page listing.\n")
			}
		}
	})
}
func Benchmark_tallyPages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// I'm not blanking the *Page values
		// before every run of tallyPages here
		// because the likelihood of
		// tallyPages calling page.cache() for
		// every page is near-zero
		if tallyPages(tallyPagesBuf); tallyPagesBuf == nil {
			b.Errorf("tallyPages() benchmark failed, got nil bytes\n")
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

type indexFields struct {
	Modtime   time.Time
	LastTally time.Time
}

var IndexCacheCases = []struct {
	name   string
	fields indexFields
	want   bool
}{
	{
		name: "test1",
		fields: indexFields{
			LastTally: time.Now(),
		},
		want: false,
	},
	{
		name: "test2",
		fields: indexFields{
			Modtime:   time.Time{},
			LastTally: time.Time{},
		},
		want: true,
	},
}

var testIndex = indexPage{
	Modtime:   time.Time{},
	LastTally: time.Time{},
}

// Check if checkCache() method on indexPage type
// is returning the expected bool
func Test_indexPage_checkCache(t *testing.T) {
	initConfigParams()
	testindexstat, err := os.Stat(confVars.assetsDir + "/" + confVars.indexFile)
	if err != nil {
		t.Errorf("Test_indexPage_checkCache(): Couldn't stat file for first test case: %v\n", err)
	}

	for _, tt := range IndexCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "test1" {
				tt.fields.Modtime = testindexstat.ModTime()
			}
			testIndex.Modtime = tt.fields.Modtime
			testIndex.LastTally = tt.fields.LastTally
			if got := testIndex.checkCache(); got != tt.want {
				t.Errorf("indexPage.checkCache() - got %v, want %v\n", got, tt.want)
			}
		})
	}
}

func Benchmark_indexPage_checkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range IndexCacheCases {
			testIndex.checkCache()
		}
	}
}

// No output to test here
func Benchmark_indexPage_cache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range IndexCacheCases {
			if err := testIndex.cache(); err != nil {
				b.Errorf("testIndex.cache() - %v\n", err)
			}
		}
	}
}

var pageCacheCase2stat, _ = os.Stat("pages/example.md")
var PageCacheCases = []struct {
	name      string
	fields    fields
	wantErr   bool
	needCache bool
}{
	{
		name: "test1",
		fields: fields{
			Longname:  "pages/test1.md",
			Shortname: "test1.md",
			Modtime:   time.Time{},
		},
		wantErr:   false,
		needCache: true,
	},
	{
		name: "example",
		fields: fields{
			Longname:  "pages/example.md",
			Shortname: "example.md",
			Modtime:   pageCacheCase2stat.ModTime(),
		},
		wantErr:   false,
		needCache: false,
	},
	{
		name: "fake page",
		fields: fields{
			Longname:  "pages/fakepage.md",
			Shortname: "fakepage.md",
			Modtime:   time.Time{},
		},
		wantErr:   true,
		needCache: false,
	},
}

// Check that the page.Body field isn't nil after
// calling page.cache(), if the page is supposed
// to exist.
func TestPage_cache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			if err := page.cache(); tt.wantErr == false {
				cachedpage := cachedPages[tt.fields.Shortname]
				if cachedpage.Body == nil {
					t.Errorf("page.cache(): got nil page body: %v\n", err)
				}
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
			if err := page.cache(); err != nil && tt.wantErr == false {
				b.Errorf("While benchmarking page.cache, caught: %v\n", err)
			}
		}
	}
}

// Make sure it's returning the appropriate
// bool for zeroed modtime and current modtime
func TestPage_checkCache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
				Modtime:   tt.fields.Modtime,
			}
			got := page.checkCache()
			if got != tt.needCache {
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
			page.checkCache()
		}
	}
}

// No output to test for
func Benchmark_genPageCache(b *testing.B) {
	initConfigParams()
	log.SetOutput(hush)
	for i := 0; i < b.N; i++ {
		genPageCache()
	}
}
