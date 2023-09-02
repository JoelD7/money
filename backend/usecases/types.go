package usecases

type authRequestBody struct {
	email    string
	password string
}

func (a authRequestBody) LogName() string {
	return "request_body"
}

func (a authRequestBody) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_email":    a.email,
		"s_password": a.password,
	}
}

type refreshTokenValue struct {
	value string
}

func (r refreshTokenValue) LogName() string { return "refresh_token" }

func (r refreshTokenValue) LogProperties() map[string]interface{} {
	return map[string]interface{}{
		"s_value": r.value,
	}
}
