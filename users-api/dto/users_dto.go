package dto

type UserDto struct {
	User_id  int64  `json:"user_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Admin    bool   `json:"admin"`
}

type UsersDto []UserDto
