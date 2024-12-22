package model

type RegisterUserResponse struct {
	Body struct {
		Status string `json:"status" example:"created" doc:"Status of creating user request"`
		Error  string `json:"error,omitempty" example:"user already exist" doc:"Errors in user creation"`
	}
}
