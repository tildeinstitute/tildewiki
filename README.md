# Tildewiki

A wiki engine designed for the needs of the [tildeverse](https://tildeverse.org)

Currently, it is in very early development. Don't try to use it yet.

## About

* Markdown rendering of all files
* Specify a file or a URL for the CSS file
* Dynamically generates index of pages and places at anchor-point in `wiki.md`
* Runs as a multithreaded service, rather than via CGI
* Easily use Nginx to proxy requests to it. This allows you to use your existing SSL certificates.
* Speed is a priority
