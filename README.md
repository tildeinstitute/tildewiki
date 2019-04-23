# Tildewiki

A wiki engine designed for the needs of the [tildeverse](https://tildeverse.org)

Uses a patched copy of [russross/blackfriday](https://github.com/russross/blackfriday) ([gopkg](https://gopkg.in/russross/blackfriday.v2)) as the markdown parser. The patch allows injection of arbitrary `<meta.../>` tags into the document header during the `markdown->html` translation.  I'll be submitting a PR of my change once I patch the development codebase.

The patched repository lives at: [gbmor-forks/blackfriday.v2-patched](https://github.com/gbmor-forks/blackfriday.v2-patched)

## About

* Markdown rendering of all files
* YAML for configuration
* Watches config file for changes and automatically reloads
* Specify a file or a URL for the CSS file
* Dynamically generates index of pages and places at anchor-point in `wiki.md`
* Runs as a multithreaded service, rather than via CGI
* Easily use Nginx to proxy requests to it. This allows you to use your existing SSL certificates.
* Speed is a priority
