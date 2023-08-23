package usecases

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/hash"
	"github.com/gbrlsnchs/jwt/v3"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

const (
	emailRegex   = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-]+$"
	passwordCost = bcrypt.DefaultCost
)

var (
	accessTokenAudience  = env.GetString("TOKEN_AUDIENCE", "https://localhost:3000")
	accessTokenIssuer    = env.GetString("TOKEN_ISSUER", "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging")
	accessTokenScope     = env.GetString("TOKEN_SCOPE", "read write")
	privateSecretName    = env.GetString("TOKEN_PRIVATE_SECRET", "staging/money/rsa/private")
	publicSecretName     = env.GetString("TOKEN_PUBLIC_SECRET", "staging/money/rsa/public")
	kidSecretName        = env.GetString("KID_SECRET", "staging/money/rsa/kid")
	accessTokenDuration  = env.GetInt("ACCESS_TOKEN_DURATION", 300)
	refreshTokenDuration = env.GetInt("REFRESH_TOKEN_DURATION", 2592000)

	errInvalidTokenLength = errors.New("invalid token length")
)

type UserCreator interface {
	CreateUser(ctx context.Context, fullName, email, password string) error
}

type UserUpdater interface {
	UpdateUser(ctx context.Context, user *models.User) error
}

type Logger interface {
	Warning(eventName string, err error, objects ...models.LoggerObject)
	Error(eventName string, err error, objects ...models.LoggerObject)
}

type InvalidTokenCache interface {
	GetInvalidTokens(ctx context.Context, email string) ([]*models.InvalidToken, error)
	AddInvalidToken(ctx context.Context, email, token string, ttl int64) error
}

type SecretManager interface {
	GetSecret(ctx context.Context, name string) (string, error)
}

// NewUserCreator creates a new user with password.
func NewUserCreator(userCreator UserCreator, logger Logger) func(ctx context.Context, fullName, email, password string) error {
	return func(ctx context.Context, fullName, email, password string) error {
		err := validateCredentials(email, password)
		if err != nil {
			logger.Error("credentials_validation_failed", err, nil)

			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
		if err != nil {
			logger.Error("password_hashing_failed", err, nil)

			return err
		}

		err = userCreator.CreateUser(ctx, fullName, email, string(hashedPassword))
		if err != nil && errors.Is(err, models.ErrExistingUser) {
			logger.Warning("user_creation_failed", err, nil)

			return err
		}

		if err != nil {
			logger.Error("sign_up_process_failed", err, nil)

			return err
		}

		return nil
	}
}

// NewUserAuthenticator authenticates a user.
func NewUserAuthenticator(userGetter UserGetter, logger Logger) func(ctx context.Context, email, password string) (*models.User, error) {
	return func(ctx context.Context, email, password string) (*models.User, error) {
		err := validateCredentials(email, password)
		if err != nil {
			logger.Error("credentials_validation_failed", err, nil)

			return nil, err
		}

		user, err := userGetter.GetUserByEmail(ctx, email)
		if err != nil {
			logger.Error("user_fetching_failed", err, nil)

			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			logger.Error("password_mismatch", err, authRequestBody{email, password})

			return nil, models.ErrWrongCredentials
		}

		return user, nil
	}
}

// NewUserTokenGenerator generates access and refresh tokens for the user.
func NewUserTokenGenerator(userUpdater UserUpdater, secretManager SecretManager, logger Logger) func(ctx context.Context, user *models.User) (*models.AuthToken, *models.AuthToken, error) {
	return func(ctx context.Context, user *models.User) (*models.AuthToken, *models.AuthToken, error) {
		now := time.Now()

		accessTokenExpiry := jwt.NumericDate(now.Add(time.Duration(accessTokenDuration) * time.Second))

		accessTokenPayload := &jwt.Payload{
			Issuer:         accessTokenIssuer,
			Subject:        user.Email,
			Audience:       jwt.Audience{accessTokenAudience},
			ExpirationTime: accessTokenExpiry,
			IssuedAt:       jwt.NumericDate(now),
		}

		accessToken, err := generateJWT(secretManager, accessTokenPayload, accessTokenScope)
		if err != nil {
			logger.Error("generate_access_token_failed", err)

			return nil, nil, err
		}

		refreshTokenExpiry := jwt.NumericDate(now.Add(time.Duration(refreshTokenDuration) * time.Second))

		refreshTokenPayload := &jwt.Payload{
			Subject:        user.Email,
			ExpirationTime: refreshTokenExpiry,
		}

		refreshToken, err := generateJWT(secretManager, refreshTokenPayload, "")
		if err != nil {
			logger.Error("generate_refresh_token_failed", err)

			return nil, nil, err
		}

		hashedAccess, err := hash.Apply(accessToken)
		if err != nil {
			logger.Error("hashing_access_token_failed", err, user)

			return nil, nil, err
		}

		hashedRefresh, err := hash.Apply(refreshToken)
		if err != nil {
			logger.Error("hashing_refresh_token_failed", err, user)

			return nil, nil, err
		}

		user.RefreshToken = hashedRefresh
		user.AccessToken = hashedAccess

		err = userUpdater.UpdateUser(ctx, user)
		if err != nil {
			logger.Error("update_user_failed", err, user)

			return nil, nil, err
		}

		access := &models.AuthToken{
			Value:      accessToken,
			Expiration: accessTokenExpiry.Time,
		}

		refresh := &models.AuthToken{
			Value:      refreshToken,
			Expiration: refreshTokenExpiry.Time,
		}

		return access, refresh, nil
	}
}

// NewRefreshTokenValidator validates a refresh token.
func NewRefreshTokenValidator(userGetter UserGetter, logger Logger) func(ctx context.Context, refreshToken string) (*models.User, error) {
	return func(ctx context.Context, refreshToken string) (*models.User, error) {
		payload, err := getTokenPayload(refreshToken)
		if err != nil {
			logger.Error("get_refresh_token_payload_failed", err)

			return nil, fmt.Errorf("%w: %v", models.ErrInvalidToken, err)
		}

		user, err := userGetter.GetUserByEmail(ctx, payload.Subject)
		if err != nil {
			logger.Error("get_user_failed", err)

			return nil, err
		}

		err = validateRefreshToken(user, refreshToken)
		if err != nil {
			logger.Warning("refresh_token_validation_failed", err, user, refreshTokenValue{refreshToken})

			return nil, fmt.Errorf("%w: %v", models.ErrInvalidToken, err)
		}

		return user, nil
	}
}

func validateCredentials(email, password string) error {
	regex := regexp.MustCompile(emailRegex)

	if email == "" {
		return models.ErrMissingEmail
	}

	if !regex.MatchString(email) {
		return models.ErrInvalidEmail
	}

	if password == "" {
		return models.ErrMissingPassword
	}

	return nil
}

func generateJWT(secrets SecretManager, payload *jwt.Payload, scope string) (string, error) {
	priv, err := getPrivateKey(secrets)
	if err != nil {
		return "", fmt.Errorf("private key fetching failed: %w", err)
	}

	var signingHash = jwt.NewRS256(jwt.RSAPrivateKey(priv))

	p := models.JWTPayload{
		Scope:   scope,
		Payload: payload,
	}

	token, err := jwt.Sign(p, signingHash)
	if err != nil {
		return "", fmt.Errorf("jwt signing failed: %w", err)
	}

	return string(token), nil
}

func getPrivateKey(secrets SecretManager) (*rsa.PrivateKey, error) {
	privateSecret, err := secrets.GetSecret(context.Background(), privateSecretName)
	if err != nil {
		return nil, err
	}

	privatePemBlock, _ := pem.Decode([]byte(privateSecret))
	if privatePemBlock == nil || !strings.Contains(privatePemBlock.Type, "PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM private block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func getTokenPayload(token string) (*models.JWTPayload, error) {
	var payload *models.JWTPayload

	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 3 {
		return nil, errInvalidTokenLength
	}

	payloadPart, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return nil, fmt.Errorf("payload decoding failed: %w", err)
	}

	err = json.Unmarshal(payloadPart, &payload)
	if err != nil {
		return nil, fmt.Errorf("paylaod unmarshalling failed: %w", err)
	}

	return payload, nil
}

func validateRefreshToken(user *models.User, refreshToken string) error {
	err := hash.CompareWithToken(user.RefreshToken, refreshToken)
	if errors.Is(err, hash.ErrHashMismatch) && user.RefreshToken != "" {
		return fmt.Errorf("%w", models.ErrRefreshTokenMismatch)
	}

	return err
}

// NewTokenInvalidator invalidates a user's tokens.
func NewTokenInvalidator(tokenCache InvalidTokenCache, logger Logger) func(ctx context.Context, user *models.User) error {
	return func(ctx context.Context, user *models.User) error {
		accessTokenTTL := time.Now().Add(time.Second * time.Duration(accessTokenDuration)).Unix()
		refreshTokenTTL := time.Now().Add(time.Second * time.Duration(refreshTokenDuration)).Unix()

		err := tokenCache.AddInvalidToken(ctx, user.Email, user.AccessToken, accessTokenTTL)
		if err != nil {
			logger.Error("access_token_invalidation_failed", err, user)

			return err
		}

		err = tokenCache.AddInvalidToken(ctx, user.Email, user.RefreshToken, refreshTokenTTL)
		if err != nil {
			logger.Error("refresh_token_invalidation_failed", err, user)

			return err
		}

		return nil
	}
}

// GetJsonWebKeySet returns a JWKS using the public and kid secret names passed in.
func GetJsonWebKeySet(secrets SecretManager, logger Logger) (*models.Jwks, error) {
	publicKey, err := getPublicKey(secrets)
	if err != nil {
		logger.Error("public_key_fetching_failed", err, nil)

		return nil, err
	}

	kid, err := getKidFromSecret(secrets)
	if err != nil {
		logger.Error("kid_fetching_failed", err, nil)

		return nil, err
	}

	return &models.Jwks{
		Keys: []models.Jwk{
			{
				Kid: kid,
				Kty: "RSA",
				Use: "sig",
				N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
				E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
			},
		},
	}, nil
}

func getPublicKey(secrets SecretManager) (*rsa.PublicKey, error) {
	publicSecret, err := secrets.GetSecret(context.Background(), publicSecretName)
	if err != nil {
		return nil, err
	}

	publicPemBlock, _ := pem.Decode([]byte(publicSecret))
	if publicPemBlock == nil || !strings.Contains(publicPemBlock.Type, "PUBLIC KEY") {
		return nil, fmt.Errorf("failed to decode PEM public block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

// The kid of the JWK that contains the public key.
// Is stored in a secret so that the lambda-authorizer can have access to it to verify that the key received is the
// right one.
func getKidFromSecret(secrets SecretManager) (string, error) {
	kidSecret, err := secrets.GetSecret(context.Background(), kidSecretName)
	if err != nil {
		return "", err
	}

	return kidSecret, nil
}

func NewUserLogout(userGetter UserGetter, tokenCache InvalidTokenCache, logger Logger) func(ctx context.Context, token string) error {
	return func(ctx context.Context, token string) error {
		payload, err := getTokenPayload(token)
		if err != nil {
			logger.Error("get_token_payload_failed", err)

			return err
		}

		user, err := userGetter.GetUserByEmail(ctx, payload.Subject)
		if err != nil {
			logger.Error("fetching_user_from_storage_failed", err, nil)

			return err
		}

		invalidateTokens := NewTokenInvalidator(tokenCache, logger)

		err = invalidateTokens(ctx, user)
		if err != nil {
			return err
		}

		return nil
	}
}
