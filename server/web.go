package server

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"strings"

	"alox.sh"
	"alox.sh/webpage"
)

type WebHandler func(*Web, http.ResponseWriter, *http.Request)

type Web struct {
	*Server
	// https://hackandsla.sh/posts/2021-11-06-serve-spa-from-go/
	// use http.FileSystem instead?
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

func (web *Web) WriteFile(responseWriter http.ResponseWriter, request *http.Request, key, name string) (err error) {
	var buffer []byte

	buffer, err = fs.ReadFile(web.FS[key], strings.TrimPrefix(name, "/"))
	if err != nil {
		return
	}

	// http.ServeContent(responseWriter, request, fileStat.Name(), fileStat.ModTime(), bytes.NewReader(buffer))

	nameByDot := strings.Split(name, ".")
	contentType := mime.TypeByExtension("." + nameByDot[len(nameByDot)-1])
	if contentType == "" {
		contentType = http.DetectContentType(buffer)
	}

	header := responseWriter.Header()
	header.Set("Content-Length", fmt.Sprintf("%d", len(buffer)))
	header.Set("Content-Type", contentType)

	responseWriter.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		_, err = responseWriter.Write(buffer)
	}

	return
}
