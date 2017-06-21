package types

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Login       string `json:"login"`
	Age         int    `json:"age"`
	Telephone   string `json:"telephone"`
	CreatedDate string `json:"created_date"`
}
