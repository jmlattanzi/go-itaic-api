package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/jmlattanzi/itaic/pc"
	"goji.io"
	"goji.io/pat"
)

func main() {
	fmt.Println("[ * ] Starting API....")

	router := goji.NewMux()

	// get all posts
	// get single post
	// add post
	router.HandleFunc(pat.Get("/post"), pc.HandleGetPosts())
	// edit post
	// delete post

	// create comment
	// delete comment
	// edit comment

	// get user
	// register user

	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router))
}
