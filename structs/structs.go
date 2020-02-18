package structs

import (
    "github.com/jinzhu/gorm"
)

type Post struct {
	Username string
	PostMessageid int
	AuthorId int 
	Text string
	Date string
	Flag int
	Image string
}


type User struct  {
	gorm.Model
	User_id int `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Username string `gorm:"not null"`
	Email string `gorm:"not null"`
	Pw_hash string `gorm:"not null"`
	Image_url string
}

type Follower struct  {
	gorm.Model
	Who_id int`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Whom_id int
}

type Message struct  {
	gorm.Model
	Message_id int`gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Author_id int `gorm:"not null"`
	Text string `gorm:"not null"`
	Pub_date int
	Flagged int
}