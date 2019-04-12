package pc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"goji.io/pat"
	"google.golang.org/api/iterator"

	"github.com/fatih/structs"
	"github.com/jmlattanzi/itaic/models"
)

// HandleGetPosts ... Gets all posts from the DB
func HandleGetPosts(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		posts := []models.Post{}
		iter := client.Collection("posts").Documents(ctx)
		for {
			post := models.Post{}
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				log.Fatal("[ ! ] Error iterating documents: ", err)
			}

			// fmt.Println(doc.Data())
			err = doc.DataTo(&post)
			if err != nil {
				log.Fatal("[ ! ] Error mapping data to struct: ", err)
			}
			posts = append(posts, post)
		}

		json.NewEncoder(res).Encode(&posts)
	}
}

// HandleGetPostByID ... Gets a single post based on id
func HandleGetPostByID(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		post := models.Post{}
		id := pat.Param(req, "id")
		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error finding document: ", err)
		}

		err = doc.DataTo(&post)
		if err != nil {
			log.Fatal("[ ! ] Error mapping data into struct: ", err)
		}

		json.NewEncoder(res).Encode(&post)
	}
}

//HandleCreatePost ...Inserts a post to the DB
func HandleCreatePost(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		// setup the new post
		newPost := models.Post{}
		err := json.NewDecoder(req.Body).Decode(&newPost)
		if err != nil {
			log.Fatal("[ ! ] Error decoding request body: ", err)
		}

		// create a new document in the collection
		doc := client.Collection("posts").NewDoc()

		// using the new doc, set the id in the post to the doc's id
		newPost.ID = doc.ID
		newPost.Created = time.Now().String()

		// write data to the doc
		_, err = doc.Create(ctx, newPost)
		if err != nil {
			log.Fatal("[ ! ] Error creating new document: ", err)
		}

		json.NewEncoder(res).Encode(&newPost)
	}
}

// HandleDeletePost ...Deletes a document form the DB
func HandleDeletePost(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		id := pat.Param(req, "id")
		_, err := client.Collection("posts").Doc(string(id)).Delete(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error deleting post: ", err)
		}

		json.NewEncoder(res).Encode("Post deleted")
	}
}

// HandleEditPost ...Edits a post in the DB
func HandleEditPost(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		type Caption struct {
			Caption string `json:"caption"`
		}

		// setup some variable
		var newCaption Caption
		id := pat.Param(req, "id")
		currentPost := models.Post{}

		// get the document
		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error getting document: ", err)
		}

		// populate the current post
		err = doc.DataTo(&currentPost)
		if err != nil {
			log.Fatal("[ ! ] Error mapping data to the struct: ", err)
		}

		// decode the new caption string
		err = json.NewDecoder(req.Body).Decode(&newCaption)
		if err != nil {
			log.Fatal("[ ! ] Error decoding request body")
		}
		currentPost.Caption = newCaption.Caption

		// this is why I have to do in such a long manner
		mappedPost := structs.Map(currentPost)
		_, err = client.Collection("posts").Doc(id).Set(ctx, mappedPost)
		if err != nil {
			log.Fatal("[ ! ] Error setting document: ", err)
		}

		json.NewEncoder(res).Encode(currentPost)
	}
}

// HandleLikePost ... Handles liking a post
func HandleLikePost(ctx context.Context, client *firestore.Client) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		id := pat.Param(req, "id")
		uid := pat.Param(req, "uid")
		user := models.User{}
		post := models.Post{}

		query := client.Collection("users").Where("uid", "==", uid)
		iter := query.Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}

			if err != nil {
				log.Fatal("[ ! ] Error iterating documents: ", err)
			}

			err = doc.DataTo(&user)
			if err != nil {
				log.Fatal("[ ! ] Error mapping data to struct: ", err)
			}
		}

		doc, err := client.Collection("posts").Doc(id).Get(ctx)
		if err != nil {
			log.Fatal("[ ! ] Error finding post")
		}

		err = doc.DataTo(&post)
		if err != nil {
			log.Fatal("[ ! ] Error mapping data to struct: ", err)
		}

		likes := user.Likes
		fmt.Println("likes: ", likes)

		addToLikes, i := remove(likes, id)
		if addToLikes == false && i == 0{
			likes = append(likes, id)
			post.Likes++
		} else {
			likes = append(likes[:i], likes[i+1:]...)
			post.Likes--
		}

		user.Likes = likes
		_, err = client.Collection("posts").Doc(id).Set(ctx, post)
		if err != nil {
			log.Fatal("[ ! ] Error setting post: ", err)
		}

		_, err = client.Collection("users").Doc(user.ID).Set(ctx, user)
		if err != nil {
			log.Fatal("[ ! ] Error setting user: ", err)
		}

		json.NewEncoder(res).Encode(&post)
	}
}

func remove(likes []string, id string) (bool, int) {
	for i := 0; i < len(likes); i++ {
		likeID := likes[i]
		if likeID == id {
			fmt.Println("id found")
			return true, i
		}
	}
	fmt.Println("id not found in list")
	return false, 0
}
