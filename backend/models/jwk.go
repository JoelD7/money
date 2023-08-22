package models

type Jwks struct {
	Keys []Jwk `json:"keys"`
}

type Jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid,omitempty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}
