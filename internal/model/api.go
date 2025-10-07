package model

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id        string `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	UserType  string `db:"user_type" json:"user_type"`
	Login     string `db:"login" json:"login"`
	Password  string `db:"password" json:"password"`
	Height    int    `db:"height" json:"height"`
	Weight    int    `db:"weight" json:"weight"`
	Age       int    `db:"age" json:"age"`
	Gender    string `db:"sex" json:"gender"`
}

type GetUserParams struct {
	LargeHeight int    `json:"large_height"` // все что больше...
	LargeWeight int    `json:"large_weight"` // все что больше...
	LargeAge    int    `json:"large_age"`    // все что больше...
	Gender      string `json:"gender"`       //
}
