package main

import (
	"github.com/russross/blackfriday"
)

func setupMarkdown(css string, title string) *blackfriday.HTMLRenderer {
	var params = blackfriday.HTMLRendererParameters{
		CSS:   css,
		Title: title,
		Flags: blackfriday.CompletePage,
	}
	return blackfriday.NewHTMLRenderer(params)
}

func render(data []byte, css string, title string) []byte {
	return blackfriday.Run(data, blackfriday.WithRenderer(setupMarkdown(css, title)))
}
