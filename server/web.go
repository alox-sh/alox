package server

import (
	"io/fs"
	"net/http"
	"strings"

	"alox.sh"
	"alox.sh/webpage"
)

type WebHandler func(*Web, http.ResponseWriter, *http.Request)

type Web struct {
	*Server
	FS       map[string]fs.FS
	Webpages map[string]*webpage.Webpage
}

func NewWeb(handler WebHandler) (web *Web) {
	web = &Web{
		Server:   NewServer(),
		FS:       map[string]fs.FS{},
		Webpages: map[string]*webpage.Webpage{},
	}
	return web.setHandler(handler)
}

func (web *Web) setHandler(handler WebHandler) *Web {
	web.Server.setHandler(func(_ alox.Server, responseWriter http.ResponseWriter, request *http.Request) {
		handler(web, responseWriter, request)
	})
	return web
}

func (web *Web) Webpage(key string) *webpage.Webpage {
	if web.Webpages[key] != nil {
		webpage, _ := web.Webpages[key].Clone()
		return webpage
	}

	return nil
}

func (web *Web) SetWebpage(key string, webpage *webpage.Webpage) *Web {
	web.Webpages[key] = webpage
	return web
}

func (web *Web) WriteFile(responseWriter http.ResponseWriter, key, name string) (err error) {
	var data []byte

	data, err = fs.ReadFile(web.FS[key], strings.TrimPrefix(name, "/"))
	if err != nil {
		return
	}

	// http.ServeContent(responseWriter, request, fileStat.Name(), fileStat.ModTime(), bytes.NewReader(buffer))

	header := responseWriter.Header()
	header.Set("Content-Type", http.DetectContentType(data))

	responseWriter.Write(data)
	return
}
