package auth

type AuthRequestBody struct {
	InitData string `json:"initData"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type TGUserData struct {
	ID int64 `json:"id"`
}
