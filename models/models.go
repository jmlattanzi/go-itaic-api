package models

// const data = {
// 	id: 'test',
// 	uid: 'test',
// 	username: 'test',
// 	caption: 'test',
// 	image: 'test',
// 	likes: 'test',
// 	created: 'test',
// 	comments: [
// 		{
// 			id: 'test',
// 			uid: 'test',
// 			comment: 'test',
// 			time: 'test',
// 			likes: 'test',
// 		},
// 	],
// }

// Post ...defines our post structure
type Post struct {
	ID       string `json:"id"`
	UID      string `json:"UID"`
	Username string `json:"username"`
	Caption  string `json:"caption"`
	ImageURL string `json:"imageURL"`
	Likes    int    `json:"likes"`
	Created  string `json:"created"`
	Comments []Comment
}

// Comment ...defines our comment structure
type Comment struct {
	ID      string `json:"comment_id"`
	UID     string `json:"comment_uid"`
	Comment string `json:"comment"`
	Created string `json:"created"`
	Likes   int    `json:"likes"`
}
