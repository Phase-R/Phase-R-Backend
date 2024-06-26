package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Age      int    `json:"age"`
	Access   string `json:"access"`
}
