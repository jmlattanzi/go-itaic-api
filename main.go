package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gorilla/handlers"
	"github.com/jmlattanzi/itaic/cc"
	"github.com/jmlattanzi/itaic/pc"
	"github.com/jmlattanzi/itaic/uc"
	"goji.io"
	"goji.io/pat"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("[ * ] Starting API....")
	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile("itaic-key.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	auth, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("[ ! ] Error getting Auth client: %v\n", err)
	}

	router := goji.NewMux()

	// post routes
	router.HandleFunc(pat.Get("/posts"), pc.HandleGetPosts(ctx, client))
	router.HandleFunc(pat.Post("/posts"), pc.HandleCreatePost(ctx, client))
	router.HandleFunc(pat.Get("/posts/:id"), pc.HandleGetPostByID(ctx, client))
	router.HandleFunc(pat.Put("/posts/:id"), pc.HandleEditPost(ctx, client))
	router.HandleFunc(pat.Delete("/posts/:id/:uid"), pc.HandleDeletePost(ctx, client))
	router.HandleFunc(pat.Put("/posts/like/:id/:uid"), pc.HandleLikePost(ctx, client))

	// comment routes
	router.HandleFunc(pat.Post("/comment/:id"), cc.HandleAddComment(ctx, client))
	router.HandleFunc(pat.Delete("/comment/:id/:comment"), cc.HandleDeleteComment(ctx, client))
	router.HandleFunc(pat.Put("/comment/:id/:comment"), cc.HandleEditComment(ctx, client))
	router.HandleFunc(pat.Put("/comment/like/:post_id/:id/:uid"), cc.HandleLikeComment(ctx, client))

	// user routes
	router.HandleFunc(pat.Get("/user/:uid"), uc.HandleGetUser(ctx, client))
	router.HandleFunc(pat.Post("/user"), uc.HandleRegisterUser(ctx, client, auth))
	router.HandleFunc(pat.Put("/user/:uid"), uc.HandleEditUser(ctx, client))

	fmt.Println("[ + ] API Started")
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router))
}
