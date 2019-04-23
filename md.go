package main

import (
	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

// sets parameters for the markdown->html renderer
func setupMarkdown(css string, title string) *bf.HTMLRenderer {
	var params = bf.HTMLRendererParameters{
		CSS:   css,
		Title: title,
		Meta: map[string]string{
			"name=\"viewport\"": "width=device-width, initial-scale=1",
		},
		Flags: bf.CompletePage | bf.Safelink,
	}
	return bf.NewHTMLRenderer(params)
}

// render the markdown into html
func render(data []byte, css string, title string) []byte {
	return bf.Run(data, bf.WithRenderer(setupMarkdown(css, title)))
}
