package models

// comment json
// {
//     "id": "369",
//     "uid": "test",
//     "username": "new post",
//     "caption": "new post",
//     "image": "test",
//     "likes": 0,
//     "created": "test",
//     "comments": [
//         {
//             "id": "test",
//             "uid": "test",
//             "comment": "test",
//             "time": "test",
//             "likes": 0
//         }
//     ]
// }

// Comment ... Defines the structure of a comment in the post
type Comment struct {
	ID      string `firestore:"id"`
	UID     string `firestore:"uid"`
	Comment string `firestore:"comment"`
	Created string `firestore:"created"`
	Likes   int    `firestore:"likes"`
}

// Post ... Defines the structure of our post in firestore
type Post struct {
	ID       string    `firestore:"id"`
	UID      string    `firestore:"uid"`
	Username string    `firestore:"username"`
	Caption  string    `firestore:"caption"`
	ImageURL string    `firestore:"imageURL"`
	Likes    int       `firestore:"likes"`
	Created  string    `firestore:"created"`
	Comments []Comment `firestore:"comments"`
}

// User ... Defines what will be stored in the user object
type User struct {
	UID        string   `firestore:"uid"`
	ID         string   `firestore:"id"`
	Username   string   `firestore:"username"`
	Email      string   `firestore:"email"`
	Bio        string   `firestore:"bio"`
	ProfilePic string   `firestore:"profile_pid"`
	Posts      []string `firestore:"posts"`
	Likes      []string `firestore:"likes"`
	Following  []string `firestore:"following"`
	Followers  []string `firestore:"followers"`
}

// {
// 	"UID": "test",
// 	"ID": "test",
// 	"Username": "test",
// 	"Email": "test@gmail.com",
// 	"Bio": "test",
//  "ProfilePic": "nice",
// 	"Posts": ["test"],
// 	"Likes": ["test"],
// 	"Following": ["test"],
// 	"Followers": ["test"]
// }
