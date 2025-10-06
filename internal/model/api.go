package model

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Uid       string `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserType  string `json:"user_type"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type GetUserParams struct {
	Uid string `json:"uid"`
}
