package dto

type LoginDtoRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDtoResponse struct {
	User_id int    `json:"user_id"`
	Token   string `json:"token"`
	Admin   bool   `json:"admin"`
}

type TokenDto struct {
	Token   string `json:"token"`
	User_id int    `json:"id_user"`
	Admin   bool   `json:"admin"`
}
