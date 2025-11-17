package entity

type User struct {
	ID   int64
	Name string
}

func (u User) CanEditComment(c Comment) bool {
	return u.ID == c.UserID
}
