package api

import "net/http"

// ShowUI redirects to the /ui/ url path
func ShowUI(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, "/ui/", 301)
}
