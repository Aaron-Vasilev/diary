package model

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	Id         int
	CreatedAt  string
	Email      string
	Name       string
	Role       Role
	SubId      string
	Subscribed bool
}

type Question struct {
	Id        int
	Text      string
	ShownDate string
}

type Note struct {
	Id          int
	UserId      int
	QuestionId  int
	Text        string
	CreatedDate string
}
