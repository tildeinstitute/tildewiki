# Tildewiki [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/tildewiki)](https://goreportcard.com/report/github.com/gbmor/tildewiki) [![GolangCI](https://img.shields.io/badge/golangci-check-blue.svg)](https://golangci.com/r/github.com/gbmor/tildewiki) [![Travis CI](https://api.travis-ci.org/gbmor/tildewiki.svg?branch=master)](https://travis-ci.org/gbmor/tildewiki)
A wiki engine designed around the needs of the [tildeverse](https://tildeverse.org)

## [v0.5.3](https://github.com/gbmor/tildewiki/releases/tag/v0.5.3)
A ton of refactoring has gone into `v0.5`. Here are some noteworthy changes:
* Various performance improvements
* Index page is now being cached
* Refresh interval for the index page is configurable
* Logging can be output to `stdout` (default), to a file, or to `/dev/null` for some peace and quiet.
* Fixed an annoying bug where a CSS change in `tildewiki.yml` wasn't reflected without a restart
* Code readability improvements

### [Development Branch](https://github.com/gbmor/tildewiki/tree/dev)
Contains all the new changes going into the next version

### Currently powering the [tilde.institute](https://tilde.institute) wiki: 
* [https://wiki.tilde.institute](https://wiki.tilde.institute) 
* [gtmetrix report](https://gtmetrix.com/reports/wiki.tilde.institute/F1tzxEch)

## Features
* Speed is a priority
* Mobile-friendly pages
* Markdown!
* Uses [kognise/water.css](https://github.com/kognise/water.css) dark theme by
default (and includes as an example, a simple but nice local CSS file)
* `YAML` configuration
* Automatically reloads config file when a change is detected.
* Generates list of pages, then places at an anchor comment in the index page
* Caches pages to memory and only re-renders when the file changes
* Very configurable. For example:
  * URL path for viewing pages
  * Directory for page data
  * File to use for index page
  * Logging output (file, `stdout`, `null`) and file location
* Runs as a multithreaded service, rather than via CGI
* Easily use `Nginx` to proxy requests to it. This allows you to use your
existing SSL certificates.

## Benchmarks

* [bombardier](https://github.com/codesenberg/bombardier)

```
bombardier -c 100 -n 200000 http://localhost:8080
Bombarding http://localhost:8080 with 200000 request(s) using 100 connection(s)
 200000 / 200000 [===========================================] 100.00% 7512/s 26s
Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      7548.57     663.04   10453.06
  Latency       13.24ms     2.38ms    49.32ms
  HTTP codes:
    1xx - 0, 2xx - 200000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:     8.55MB/s

```

* [baton](https://github.com/americanexpress/baton)
```
$ baton -u http://localhost:8080 -c 100 -r 200000

...

=========================== Results ========================================

Total requests:                                200000
Time taken to complete requests:        27.270626274s
Requests per second:                             7334
Max response time (ms):                            52
Min response time (ms):                             0
Avg response time (ms):                         13.11

========= Percentage of responses by status code ==========================

Number of connection errors:                        0
Number of 1xx responses:                            0
Number of 2xx responses:                       200000
Number of 3xx responses:                            0
Number of 4xx responses:                            0
Number of 5xx responses:                            0

========= Percentage of responses received within a certain time (ms)======

         9% : 5 ms
        13% : 10 ms
        79% : 15 ms
        95% : 20 ms
        98% : 25 ms
        99% : 30 ms
        99% : 35 ms
        99% : 40 ms
        99% : 45 ms
       100% : 52 ms

===========================================================================

```

### Notes
* Builds with `Go 1.11` and `Go 1.12`. Not tested with any other version.
* Tested on Linux (Ubuntu 18.04LTS, Debian 9) and OpenBSD 6.4
* If you have access to other environments and can test, please let me know.
It will be much appreciated.

For [tildeverse](https://tildeverse.org) projects, we tend to use a PR
workflow. For example, wiki pages are submitted to the repo via pull
request. That's what I'm initially designing this around. I will likely
add authentication and in-place page editing last, after everything else
is done, including unit tests.

Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday)
([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown
parser. The patch allows injection of various `<meta.../>` tags into
the document header during the `markdown->html` translation.

* The patched `v2` repository lives at:
[gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)
* The patched `master` repo lives at:
[gbmor-forks/blackfriday](https://github.com/gbmor-forks/blackfriday).
* The PR can be found here: [allow writing of user-specified
&lt;meta.../&gt;...](https://github.com/russross/blackfriday/pull/541)

