package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"

	"alox.sh/website"
)

const HTTPPort = 2000

func main() {
	website, err := website.NewDirSPA("./website", "index.html")
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Listening on :%d\n", HTTPPort)
	http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), website)

	return

	// website.ServeHTTP(nil, nil)

	//

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/foo", func(responseWriter http.ResponseWriter, request *http.Request) {
		page, err := website.NewPage(request)
		if err != nil {
			log.Panic(err)
		}

		page.Title = "The foo page"
		page.Meta = append(page.Meta, []html.Attribute{{
			Key: "description",
			Val: "This meta tag has been injected on the fly",
		}})

		err = page.WriteToResponse(responseWriter)
		if err != nil {
			fmt.Printf("Unexpacted error: %+v\n", err)
		}
	})

	serveMux.Handle("/", http.FileServer(http.FS(website.GetRootFS())))

	fmt.Printf("Listening on :%d\n", HTTPPort)
	http.ListenAndServe(fmt.Sprintf(":%d", HTTPPort), serveMux)
}
