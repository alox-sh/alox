package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"alox.sh/webpage"
)

type SPAParams struct {
	RootFS         fs.FS
	InputHTMLName  string
	HydrateWebpage func(*SPA, http.ResponseWriter, *http.Request, *webpage.Webpage)
}

type SPA struct {
	*Web
	SPAParams
}

func NewSPA(params SPAParams) (spa *SPA, err error) {
	var (
		inputHTMLFile     fs.File
		inputHTMLFileStat fs.FileInfo
	)

	spa = &SPA{
		SPAParams: params,
		Web: NewWeb(func(_ *Web, responseWriter http.ResponseWriter, request *http.Request) {
			spa.serve(responseWriter, request)
		}),
	}

	if len(spa.InputHTMLName) < 1 {
		spa.InputHTMLName = "index.html"
	}

	spa.FS["root"] = params.RootFS

	if inputHTMLFile, err = spa.FS["root"].Open(spa.InputHTMLName); err != nil {
		return
	}

	if inputHTMLFileStat, err = inputHTMLFile.Stat(); err != nil {
		return
	}

	if inputHTMLFileStat.IsDir() {
		err = fmt.Errorf("Input HTML file '%s' is a directory", spa.InputHTMLName)
		return
	}

	if spa.Webpages[spa.InputHTMLName], err = webpage.NewWebpage(inputHTMLFile); err != nil {
		return
	}

	return
}

func (spa *SPA) serve(responseWriter http.ResponseWriter, request *http.Request) {
	webpage := spa.Webpage(spa.InputHTMLName)

	if spa.FS["root"] == nil || webpage == nil {
		spa.HandleError(responseWriter, request, fmt.Errorf("Invalid SPA"))
		return
	}

	if err := spa.WriteFile(responseWriter, "root", request.URL.Path); err == nil {
		return
	}

	// file, err := spa.FS["root"].Open(request.URL.Path)
	// if err == nil {
	// 	fileStat, err := file.Stat()
	// 	if err == nil && !fileStat.IsDir() {
	// 		buffer, err := ioutil.ReadAll(file) // suboptimal (large files)
	// 		if err != nil {
	// 			spa.HandleError(responseWriter, request, err)
	// 			return
	// 		}

	// 		// var data []byte
	// 		// data, err = fs.ReadFile(spa.FS["root"], request.URL.Path)
	// 		// responseWriter.Write(data)

	// 		http.ServeContent(responseWriter, request, fileStat.Name(), fileStat.ModTime(), bytes.NewReader(buffer))
	// 		return
	// 	}
	// }

	if spa.HydrateWebpage != nil {
		spa.HydrateWebpage(spa, responseWriter, request, webpage)
	}

	spa.WriteWebpage(responseWriter, request, webpage)
}
