package webpage

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/html"
)

func (webpage *Webpage) addJavaScriptVariable(name, value string) {
	webpage.Scripts = append(webpage.Scripts, Script{
		Head: true,
		Content: fmt.Sprintf(`
window.%s=%s;
		`, name, value),
	})
	return
}

func (webpage *Webpage) AddJavaScriptVariable(name string, value interface{}) (err error) {
	valueJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	webpage.addJavaScriptVariable(name, string(valueJson))
	return
}

func (webpage *Webpage) AddJavaScriptRedirect(redirectURI string) {
	webpage.Scripts = append(webpage.Scripts, Script{
		Head: true,
		Content: fmt.Sprintf(`
if (window && window.location && typeof window.location.assign === "function") {
    window.location.assign("%s");
}
		`, redirectURI),
	})
}

func (webpage *Webpage) AddGA(gaID string) {
	webpage.Scripts = append(webpage.Scripts, Script{
		Src:        fmt.Sprintf("https://www.googletagmanager.com/gtag/js?id=%s", gaID),
		Attributes: []html.Attribute{{Key: "async"}},
	})

	webpage.Scripts = append(webpage.Scripts, Script{
		Head: true,
		Content: fmt.Sprintf(`
window.dataLayer = window.dataLayer || [];

function gtag() {
    dataLayer.push(arguments);
}

gtag('js', new Date());
gtag('config', '%s');
		`, gaID),
	})
}

func (webpage *Webpage) AddGTM(gtmID string) {
	webpage.Scripts = append(webpage.Scripts, Script{
		Head: true,
		Content: fmt.Sprintf(`
window.dataLayer = window.dataLayer || [];

(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','%s');
		`, gtmID),
	})

	webpage.noscriptNode.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "iframe",
		Attr: []html.Attribute{{
			Key: "src",
			Val: fmt.Sprintf("https://www.googletagmanager.com/ns.html?id=%s", gtmID),
		}, {
			Key: "height",
			Val: "0",
		}, {
			Key: "width",
			Val: "0",
		}, {
			Key: "style",
			Val: "display:none;visibility:hidden;",
		}},
	})
	return
}
