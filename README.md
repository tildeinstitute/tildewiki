# Tildewiki

A minimalist wiki engine designed for the needs of the [tildeverse](https://tildeverse.org)

Currently, it is in very early development. Don't try to use it yet.

## About

* Markdown rendering of all files
* Specify a file or a URL for the CSS file
* Dynamically generates index of pages and places at anchor-point in `wiki.md`
* Runs as a service, rather than via CGI
* Because it runs as a service, easily use Nginx to forward requests to it. This allows you to use your existing SSL certificates.
* Speed is a priority
