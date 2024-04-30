package types

type User struct {
	ID    int64  `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
