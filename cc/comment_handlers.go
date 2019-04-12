package cc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	shortid "github.com/jasonsoft/go-short-id"

	"github.com/jmlattanzi/itaic/models"
	"goji.io/pat"

	"cloud.google.com/go/firestore"
)

// HandleAddComment ... Adds a comment to the db
func HandleAddComment(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		id := pat.Param(req, "id")
		type NewComment struct {
			Comment string `json:"comment"`
		}

		opt := shortid.Options{
			Number:        14,
			StartWithYear: false,
			EndWithHost:   false,
		}

		newComment := models.Comment{
			ID:      shortid.Generate(opt),
			Created: time.Now().String(),
			Likes:   0,
			Comment: "",
			UID:     "test uid",
		}
		currentPost := models.Post{}
		err := json.NewDecoder(req.Body).Decode(&newComment)
		if err != nil {
			log.Fatal("[ ! ] Error decoding request body: ", err)
		}

		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error getting document: ", err)
		}
		err = doc.DataTo(&currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error writing data to struct: ", err)
		}

		currentPost.Comments = append(currentPost.Comments, newComment)
		_, err = client.Collection("posts").Doc(id).Set(ctx, currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error setting document: ", err)
		}

		json.NewEncoder(res).Encode(currentPost)
	}
}

// HandleDeleteComment ... Deletes a comment based on post id and comment id
func HandleDeleteComment(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		id := pat.Param(req, "id")
		commentID := pat.Param(req, "comment")
		currentPost := models.Post{}
		comments := []models.Comment{}

		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error getting document")
		}

		err = doc.DataTo(&currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error mapping data into struct: ", err)
		}

		for _, comment := range currentPost.Comments {
			if comment.ID == commentID {
				fmt.Println("[ + ] Comment found")
			} else {
				comments = append(comments, comment)
			}
		}

		currentPost.Comments = comments
		_, err = client.Collection("posts").Doc(id).Set(ctx, currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error setting document: ", err)
		}
		json.NewEncoder(res).Encode(currentPost)
	}
}

// HandleEditComment ... Edits a comment and submits to the db
func HandleEditComment(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		id := pat.Param(req, "id")
		commentID := pat.Param(req, "comment")
		currentPost := models.Post{}
		comments := []models.Comment{}

		type NewComment struct {
			Comment string `json:"comment"`
		}
		newComment := NewComment{}

		err := json.NewDecoder(req.Body).Decode(&newComment)
		if err != nil {
			log.Fatal("[ ! ] Error decoding the response body: ", err)
		}

		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error getting document")
		}

		err = doc.DataTo(&currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error mapping data into struct: ", err)
		}

		for _, comment := range currentPost.Comments {
			if comment.ID == commentID {
				fmt.Println("[ + ] Comment found")
				comment.Comment = newComment.Comment
			}

			comments = append(comments, comment)
		}

		currentPost.Comments = comments
		_, err = client.Collection("posts").Doc(id).Set(ctx, currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error setting document: ", err)
		}
		json.NewEncoder(res).Encode(currentPost)
	}
}
