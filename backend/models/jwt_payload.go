package models

import "github.com/gbrlsnchs/jwt/v3"

type JWTPayload struct {
	Scope string `json:"scope,omitempty"`
	*jwt.Payload
}
