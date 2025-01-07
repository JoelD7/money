package usecases

type authRequestBody struct {
	username string
	password string
}

func (a authRequestBody) Key() string {
	return "request_body"
}

func (a authRequestBody) Value() map[string]interface{} {
	return map[string]interface{}{
		"s_username": a.username,
		"s_password": a.password,
	}
}

type refreshTokenValue struct {
	value string
}

func (r refreshTokenValue) Key() string { return "refresh_token" }

func (r refreshTokenValue) Value() map[string]interface{} {
	return map[string]interface{}{
		"s_value": r.value,
	}
}
