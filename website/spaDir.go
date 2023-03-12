package website

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"alox.sh"
)

type DirSPA struct {
	website
	RootFS    fs.FS
	InputHTML io.Reader
}

func NewDirSPA(rootPath string, inputHTMLName string) (website alox.Website, err error) {
	var (
		rootFS            = os.DirFS(rootPath)
		inputHTMLFile     fs.File
		inputHTMLFileStat fs.FileInfo
	)

	if inputHTMLFile, err = rootFS.Open(inputHTMLName); err != nil {
		return
	}

	if inputHTMLFileStat, err = inputHTMLFile.Stat(); err != nil {
		return
	}

	if inputHTMLFileStat.IsDir() {
		err = fmt.Errorf("Input HTML file '%s' is a directory", inputHTMLName)
		return
	}

	return &DirSPA{
		website:   newWebsite(),
		RootFS:    rootFS,
		InputHTML: inputHTMLFile,
	}, nil
}

func (dirSPA *DirSPA) GetRootFS() fs.FS {
	return dirSPA.RootFS
}

func (dirSPA *DirSPA) NewPage(request *http.Request) (*alox.Page, error) {
	return alox.NewPage(dirSPA, request, dirSPA.InputHTML)
}

func (dirSPA *DirSPA) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// dirSPA.website.ServeHTTP(responseWriter, request)

	var (
		err  error
		page *alox.Page
	)

	if page, err = dirSPA.NewPage(request); err != nil {
		return
	}

	page.WriteToResponse(responseWriter)
}
