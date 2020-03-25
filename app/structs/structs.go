package structs

type Post struct {
	Username   string
	Message_id int
	Author_id  int
	Text       string
	Pub_date   string
	Flagged    int
	Image_url  string
}

type User struct {
	User_id   int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
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
	Message_id int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Author_id  int    `gorm:"not null"`
	Text       string `gorm:"not null"`
	Pub_date   string
	Flagged    int
}
