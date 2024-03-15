package model

type User struct {
  Id int
  CreatedAt string
  Email string
  Name string
  Role string
  SubId string
  Subscribed bool
}

type Question struct {
  Id int
  Text string
  ShownDate string
}

type Note struct {
  Id int
  UserId int
  QuestionId int
  Text string
  CreatedDate string
}
