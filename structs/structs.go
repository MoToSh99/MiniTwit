package structs

type Post struct {
	Username      string
	PostMessageid int
	AuthorId      int
	Text          string
	Date          string
	Flag          int
	Image         string
}

type User struct {
	User_id   int    `gorm:"PRIMARY_KEY"`
	Username  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Pw_hash   string `gorm:"not null"`
	Image_url string
}

type Follower struct {
	Who_id  int 
	Whom_id int
}

type Message struct {
	Message_id int    `gorm:"PRIMARY_KEY"`
	Author_id  int    `gorm:"not null"`
	Text       string `gorm:"not null"`
	Pub_date   string
	Flagged    int
}
