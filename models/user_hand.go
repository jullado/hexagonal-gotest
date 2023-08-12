package models

type HandRegisterBodyModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HandLoginBodyModel struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HandResetPasswordParamsModel struct {
	UserId string `params:"user_id"`
}

type HandResetPasswordBodyModel struct {
	Password string `json:"password"`
}

type HandDeleteUserParamsModel struct {
	UserId string `params:"user_id"`
}
