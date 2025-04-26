package usecases

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gbrlsnchs/jwt/v3"

	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/hash"
	"github.com/JoelD7/money/backend/shared/logger"
)

var (
	invalidJWTErrs = []error{jwt.ErrAudValidation, jwt.ErrExpValidation, jwt.ErrIatValidation, jwt.ErrIssValidation,
		jwt.ErrJtiValidation, jwt.ErrNbfValidation, jwt.ErrSubValidation}
)

// NewTokenVerifier validates a JWT against the authentication server. Returns the subject of the token if successful.
func NewTokenVerifier(jwksGetter JWKSGetter, secretManager SecretManager, tokenCache InvalidTokenCache) func(ctx context.Context, token string) (string, error) {
	return func(ctx context.Context, token string) (string, error) {
		payload, err := getTokenPayload(token)
		if err != nil {
			logger.Error("getting_token_payload_failed", err, nil)

			return "", fmt.Errorf("%v: %w", err, models.ErrUnauthorized)
		}

		response, err := jwksGetter.Get(payload.Issuer + "/auth/jwks")
		if err != nil {
			logger.Error("getting_jwks_failed", err, nil)

			return "", err
		}

		defer func() {
			closeErr := response.Body.Close()
			if closeErr != nil {
				logger.Error("closing_response_body_failed", closeErr, nil)

				err = closeErr
			}
		}()

		jwksVal := new(models.Jwks)
		err = json.NewDecoder(response.Body).Decode(jwksVal)
		if err != nil {
			logger.Error("decoding_response_body_failed", err, nil)

			return "", err
		}

		publicKey, err := getPublicKeyFromJWKS(ctx, jwksVal, secretManager)
		if err != nil {
			logger.Error("getting_public_key_failed", err, nil)

			return "", err
		}

		decryptingHash := jwt.NewRS256(jwt.RSAPublicKey(publicKey))
		receivedPayload := &jwt.Payload{}

		err = validateJWTPayload(token, receivedPayload, decryptingHash)
		if err != nil {
			logger.Error("jwt_validation_failed", err, nil)

			return "", err
		}

		err = compareAccessTokenAgainstBlacklistRedis(ctx, tokenCache, payload.Subject, token)
		if errors.Is(err, models.ErrInvalidToken) {
			logger.Warning("blacklisted_token_use_detected", err, models.Any("token", map[string]interface{}{"s_value": token}))
		}

		if err != nil {
			return "", err
		}

		return payload.Subject, nil
	}
}

func getPublicKeyFromJWKS(ctx context.Context, jwksVal *models.Jwks, secrets SecretManager) (*rsa.PublicKey, error) {
	kid, err := getKidFromSecret(ctx, secrets)
	if err != nil {
		return nil, err
	}

	var signingKey *models.Jwk

	for _, key := range jwksVal.Keys {
		if key.Kid == kid {
			signingKey = &key
		}
	}

	if signingKey == nil {
		return nil, models.ErrSigningKeyNotFound
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(signingKey.N)
	if err != nil {
		return nil, err
	}

	n := new(big.Int)
	n.SetBytes(nBytes)

	eBytes, err := base64.RawURLEncoding.DecodeString(signingKey.E)
	if err != nil {
		return nil, err
	}

	e := new(big.Int)
	e.SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func validateJWTPayload(token string, payload *jwt.Payload, decryptingHash *jwt.RSASHA) error {
	now := time.Now()

	jwtAudience := env.GetString("TOKEN_AUDIENCE", "")
	jwtIssuer := env.GetString("TOKEN_ISSUER", "")

	expValidator := jwt.ExpirationTimeValidator(now)
	issValidator := jwt.IssuerValidator(jwtIssuer)
	audValidator := jwt.AudienceValidator(jwt.Audience{jwtAudience})

	validatePayload := jwt.ValidatePayload(payload, issValidator, audValidator, expValidator)

	_, err := jwt.Verify([]byte(token), decryptingHash, payload, validatePayload)
	if isErrorInvalidJWT(err) {
		return fmt.Errorf("%v: %w", err, models.ErrInvalidToken)
	}

	if err != nil {
		return fmt.Errorf("%v: %w", err, models.ErrUnauthorized)
	}

	return nil
}

func isErrorInvalidJWT(err error) bool {
	for _, e := range invalidJWTErrs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func compareAccessTokenAgainstBlacklistRedis(ctx context.Context, tokenCache InvalidTokenCache, username, token string) error {
	invalidTokens, err := tokenCache.GetInvalidTokens(ctx, username)
	if err != nil && !errors.Is(err, models.ErrInvalidTokensNotFound) {
		return err
	}

	for _, it := range invalidTokens {
		err = hash.CompareWithToken(it.Token, token)
		if err == nil {
			logger.Warning("invalid_token_use_detected", models.ErrInvalidToken, models.Any("token_comparison", map[string]interface{}{
				"s_bearer_token":             token,
				"s_saved_invalid_token_hash": it.Token,
			}))

			return models.ErrInvalidToken
		}

		if !errors.Is(err, hash.ErrHashMismatch) {
			logger.Error("token_comparison_against_blacklist_failed", err, models.Any("token_comparison", map[string]interface{}{
				"s_bearer_token":             token,
				"s_saved_invalid_token_hash": it.Token,
			}))
		}
	}

	return nil
}
