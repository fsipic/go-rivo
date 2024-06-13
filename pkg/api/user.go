package api

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User
var nextUserID int = 1

func CreateUser(name, email, password string) User {
	user := User{
		ID:       nextUserID,
		Name:     name,
		Email:    email,
		Password: password,
	}
	users = append(users, user)
	nextUserID++
	return user
}
