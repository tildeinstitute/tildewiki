# TildeWiki [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/tildewiki)](https://goreportcard.com/report/github.com/gbmor/tildewiki) [![GolangCI](https://img.shields.io/badge/golangci-check-blue.svg)](https://golangci.com/r/github.com/gbmor/tildewiki) [![Travis CI](https://api.travis-ci.org/gbmor/tildewiki.svg?branch=master)](https://travis-ci.org/gbmor/tildewiki)
TildeWiki is a memory-caching static site server. The possible uses of TildeWiki range from blogs to wikis, and more.
Let me know if you adapt it to a new use-case, I'm always interested!

Originally designed around the needs of the [tildeverse](https://tildeverse.org).<sup><a href="#1">1</a></sup>

[\[Features\]](#features) | [\[Installation\]](#installation) | [\[Benchmarks\]](#benchmarks) | [\[Notes\]](#notes)

## [v0.6.1](https://github.com/gbmor/tildewiki/releases/tag/v0.6.1)
A ton of refactoring has gone into `v0.6`
* Various performance improvements
* Code readability improvements
* Improved testing (~61% coverage)
* Script to automate build/install
* Startup script to daemonize the process

### [Development Branch](https://github.com/gbmor/tildewiki/tree/dev)
Contains the changes going into the next version

### Currently powering the [tilde.institute](https://tilde.institute) wiki: 
* [https://wiki.tilde.institute](https://wiki.tilde.institute) 
* [gtmetrix report](https://gtmetrix.com/reports/wiki.tilde.institute/F1tzxEch)

## <a name="features"></a>Features
* Speed is a priority
* Mobile-friendly pages
* Markdown!<sup><a href="#2">2</a></sup>
* Uses [kognise/water.css](https://github.com/kognise/water.css) dark theme by
default (and includes as an example, a simple but nice local CSS file)<sup><a href="#3">3</a></sup>
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

## <a name="installation"></a>Installation

The installation script uses `bash`, and the startup script uses `daemonize`. Both should
be available in any Linux distribution's package repositories. However, they are not
required to use TildeWiki.

There's an archive containing just the binary, config file, and required directories/files
available. If you download this, you can skip to the section on setting up TildeWiki.

[download release binary](https://github.com/gbmor/tildewiki/releases/download/v0.6.0/tildewiki-0.6.1-bin-linux-amd64.tar.gz)

### Using the scripts

First, clone the repository or download and untar a release archive, then enter the directory.

```
$ git clone git://github.com/gbmor/tildewiki.git && cd tildewiki

$ curl -L https://github.com/gbmor/tildewiki/archive/v0.6.0.tar.gz | tar xzvf - && cd tildewiki-v0.6.0
```

If you used `git`, the master branch will be the most recent release. Development work stays
in the `dev` branch, so there's no need to look for a tag.

Execute `setup.sh` as root, with the `install` argument:

```
$ sudo ./setup.sh install
```

Once you receive the confirmation message, and no errors have appeared, you may run the
startup script as root to test the installation:

```
$ sudo tildewiki
```

TildeWiki will drop privileges to the `tildewiki` user, which was created by the script.

I'm going to add a `systemd` service file soon. For now, it'll need to be started like this.

### Building manually

If you prefer, you can install it this way. Clone the repository or download a source archive
like above, and enter the directory. Once in the directory, you'll need to build the binary.

```
$ go build
```

It won't take long. Also, for those new to `go`, it doesn't need to live in your `GOPATH` as
it's been set up to use Go Modules. Vendored dependencies are also included.

After it finishes, you can leave the binary where it is or move it somewhere else. Remember
to move the `pages` and `assets` directories with it, along with `tildewiki.yaml`.

### Setting up TildeWIki

Begin by combing through `tildewiki.yaml` (if you used the scripts, it's in `/usr/local/tildewiki`)
and changing the options to something appropriate to your site. Afterwards, place your markdown-formatted
pages into the directory specified by `PageDir` in the config and place your markdown-formatted 
index file, containing the anchor comment `<!--pagelist-->`, into the `AssetsDir`. Feel free to
change the favicon and CSS to your liking.

Once that's all done, either run `/usr/local/bin/tildewiki` (if you've used the scripts) or run
the binary manually.

### Serving TildeWiki

Unless you plan on serving directly from :8080 (which is fine!), or whichever port you chose in 
`tildewiki.yaml`, I recommend proxying requests to TildeWiki so it can be served from a subdomain,
for example. There are several options for this, namely [Caddy](https://caddyserver.com/) and 
[nginx](https://nginx.org). The best option is for you to use Caddy: it integrates TLS certificate 
renewal and has a *very* easy configuration syntax.

If you're going to use Nginx, here's an example server block for you to start with. Note: this 
example uses TLS and http2. [LetsEncrypt](https://letsencrypt.org) is awesome, and free. 
Their `certbot` tool is really easy to use.

```
server {
    server_name wiki.example.com;
    listen [::]:443 ssl http2;
    listen 0.0.0.0:443 ssl http2;
    ssl_certificate /etc/letsencrypt/live/wiki.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wiki.example.com/privkey.pem;
    include /etc/letsencrypt/options-ssl-nginx.conf;
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem;
    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://127.0.0.1:8080;
    }
}
server {
    if ($host = wiki.example.com) {
        return 301 https://$host$request_uri;
    }
    listen 80;
    server_name wiki.example.com;
    return 404;
}
```

## <a name="benchmarks"></a>Benchmarks

* [bombardier](https://github.com/codesenberg/bombardier)

```
$ bombardier -c 100 -n 200000 http://localhost:8080

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

## <a name="notes"></a>Notes
* Builds with `Go 1.11` and `Go 1.12`. Not tested with any other version.
* Tested on Linux (Ubuntu 18.04LTS, Debian 9) and OpenBSD 6.4

1. <a name="1"></a>For [tildeverse](https://tildeverse.org) projects, we tend to use a PR
workflow for collaboration. For example, wiki pages are submitted to the repo via pull
request. I'm currently evaluating other options for page creation and editing.

2. <a name="2"></a>Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday)
([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown
parser. The patch allows injection of various `<meta.../>` tags into
the document header during the `markdown->html` translation.

   * The patched `v2` repository lives at:
[gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)
   
   * The patched `master` repo lives at:
[gbmor-forks/blackfriday](https://github.com/gbmor-forks/blackfriday).
   
   * The PR can be found here: [allow writing of user-specified
&lt;meta.../&gt;...](https://github.com/russross/blackfriday/pull/541)

3. <a name="3"></a>The local CSS provided is the "58 bytes of CSS" from [https://jrl.ninja/etc/1/](https://jrl.ninja/etc/1/)

