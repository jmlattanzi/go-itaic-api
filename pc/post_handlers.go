package pc

import (
	"net/http"
)

func HandleGetPosts() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
	}
}
