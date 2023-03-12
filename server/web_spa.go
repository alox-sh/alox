package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"alox.sh/webpage"
)

type SPAHandler func(*SPA, http.ResponseWriter, *http.Request)

type SPAParams struct {
	FS             fs.FS
	FSKey          string
	InputHTMLName  string
	HydrateWebpage func(*SPA, http.ResponseWriter, *http.Request, *webpage.Webpage)
}

type SPA struct {
	*Web
	SPAParams
}

func NewSPA(handler SPAHandler, params SPAParams) (spa *SPA, err error) {
	var (
		inputHTMLFile     fs.File
		inputHTMLFileStat fs.FileInfo
	)

	spa = &SPA{
		SPAParams: params,
		Web: NewWeb(func(_ *Web, responseWriter http.ResponseWriter, request *http.Request) {
			if handler == nil {
				spa.WriteFileOrWebpage(responseWriter, request)
				return
			}

			handler(spa, responseWriter, request)
		}),
	}

	if len(spa.FSKey) < 1 {
		spa.FSKey = "root"
	}

	if len(spa.InputHTMLName) < 1 {
		spa.InputHTMLName = "index.html"
	}

	spa.Web.FS[spa.FSKey] = params.FS

	if inputHTMLFile, err = params.FS.Open(spa.InputHTMLName); err != nil {
		return
	}

	if inputHTMLFileStat, err = inputHTMLFile.Stat(); err != nil {
		return
	}

	if !inputHTMLFileStat.Mode().IsRegular() {
		err = fmt.Errorf("input HTML file '%s' is irregular file", spa.InputHTMLName)
		return
	}

	if spa.Webpages[spa.InputHTMLName], err = webpage.NewWebpage(inputHTMLFile); err != nil {
		return
	}

	return
}

func (spa *SPA) WriteFile(responseWriter http.ResponseWriter, request *http.Request) (err error) {
	if spa.Web.FS[spa.FSKey] == nil {
		return fmt.Errorf("invalid FS '%s'", spa.FSKey)
	}

	return spa.Web.WriteFile(responseWriter, spa.FSKey, request.URL.Path)
}

func (spa *SPA) MustWriteFile(responseWriter http.ResponseWriter, request *http.Request) {
	if err := spa.WriteFile(responseWriter, request); err != nil {
		spa.HandleError(responseWriter, request, err)
	}
}

func (spa *SPA) WriteWebpage(responseWriter http.ResponseWriter, request *http.Request) (err error) {
	webpage := spa.Webpage(spa.InputHTMLName)

	if webpage == nil {
		return fmt.Errorf("invalid webpage '%s'", spa.InputHTMLName)
	}

	return spa.Web.WriteWebpage(responseWriter, request, webpage)
}

func (spa *SPA) MustWriteWebpage(responseWriter http.ResponseWriter, request *http.Request) {
	if err := spa.WriteWebpage(responseWriter, request); err != nil {
		spa.HandleError(responseWriter, request, err)
	}
}

func (spa *SPA) WriteFileOrWebpage(responseWriter http.ResponseWriter, request *http.Request) {
	if err := spa.WriteFile(responseWriter, request); err == nil {
		return
	}

	spa.MustWriteWebpage(responseWriter, request)
}
