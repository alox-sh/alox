package main

import (
	"embed"

	"alox.sh/website"
)

//go:embed website/index.html
var indexHTML []byte

//go:embed website
var rootFS embed.FS

func main() {
	website := website.NewEmbededSPA(indexHTML, rootFS)

	_ = website
}
