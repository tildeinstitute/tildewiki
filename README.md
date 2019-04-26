# Tildewiki [![gpl-v3](https://img.shields.io/badge/license-GPLv3-brightgreen.svg "GPL v3")](https://github.com/gbmor/tildewiki/blob/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/tildewiki)](https://goreportcard.com/report/github.com/gbmor/tildewiki) [![GolangCI](https://img.shields.io/badge/golang%20ci-success-blue.svg)](https://golangci.com/r/github.com/gbmor/tildewiki)

A wiki engine designed around the needs of the [tildeverse](https://tildeverse.org)

This project is still in a pre-release stage and is not yet feature complete. Code in the repo builds and works but is subject to refactoring at any time. It will stabilize soon.

Currently powering the [tilde.institute](https://tilde.institute) wiki: 
* [https://wiki.tilde.institute](https://wiki.tilde.institute)
* [gtmetrix report](https://gtmetrix.com/reports/wiki.tilde.institute/15byZymY)

## Features

* Speed is a priority
* Mobile-friendly pages
* Markdown!
* Uses [kognise/water.css](https://github.com/kognise/water.css) dark theme by default (and includes as an example, a simple but nice local CSS file)
* Automatically reloads YAML configuration when a change is detected.
* Generates list of pages and places at anchor-point in index page
* Caches pages to memory and only re-renders when the file modification time changes
* Extremely configurable: URL path for viewing pages, directory for page data, file to use for index page, etc.
* Runs as a multithreaded service, rather than via CGI
* Easily use Nginx to proxy requests to it. This allows you to use your existing SSL certificates.

## Notes

For [tildeverse](https://tildeverse.org) projects, we tend to use a PR workflow. For example, wiki pages are submitted to the repo via pull request. That's what I'm initially designing this around. I will likely add authentication and in-place page editing last, after everything else is done.

Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday) ([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown parser. The patch allows injection of various `<meta.../>` tags into the document header during the `markdown->html` translation.

* The patched `v2` repository lives at: [gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)
* The patched `master` repo lives at: [gbmor-forks/blackfriday](https://github.com/gbmor-forks/blackfriday). 
* The PR can be found here: [allow writing of user-specified &lt;meta.../&gt;...](https://github.com/russross/blackfriday/pull/541)

