package website

import (
	"encoding/json"
	"fmt"
	"net/http"

	"alox.sh"
)

type website struct {
	disableHTTPRedirect bool
	stateGlobalName     string
}

func newWebsite() website {
	return website{stateGlobalName: "__websiteState"}
}

func (website website) GetDisableHTTPRedirect() bool {
	return website.disableHTTPRedirect
}

func (website website) SetDisableHTTPRedirect(disableHTTPRedirect bool) {
	website.disableHTTPRedirect = disableHTTPRedirect
}

func (website website) GetStateGlobalName() string {
	return website.stateGlobalName
}

func (website website) SetStateGlobalName(name string) {
	website.stateGlobalName = name
}

func (website website) MarshalStateMap(page *alox.Page) map[string]interface{} {
	return map[string]interface{}{
		"page": page.MarshalMap(),
	}
}

func (website website) MarshalStateJS(page *alox.Page) (js string, err error) {
	var (
		stateJSON    []byte
		websiteState = website.MarshalStateMap(page)
	)

	if stateJSON, err = json.Marshal(websiteState); err != nil {
		return
	}

	return fmt.Sprintf("window.%s=%s;", website.stateGlobalName, stateJSON), nil
}

func (website website) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
}
