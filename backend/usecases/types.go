package usecases

type authRequestBody struct {
	Username string
	Password string
}

func (a authRequestBody) Key() string {
	return "request_body"
}

func (a authRequestBody) Value() map[string]interface{} {
	return map[string]interface{}{
		"s_username": a.Username,
		"s_password": a.Password,
	}
}
