package main

import (
	"context"
	"fmt"
	"github.com/JoelD7/money/backend/shared/apigateway"
	"github.com/JoelD7/money/backend/shared/env"
	"github.com/JoelD7/money/backend/shared/logger"
	"github.com/JoelD7/money/backend/shared/secrets"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

var (
	privateSecretName string
	publicSecretName  string
	kidSecretName     string
)

func TestMain(m *testing.M) {
	err := env.LoadEnvTesting()
	if err != nil {
		panic(fmt.Errorf("loading environment failed: %v", err))
	}

	privateSecretName = env.GetString("TOKEN_PRIVATE_SECRET", "")
	publicSecretName = env.GetString("TOKEN_PUBLIC_SECRET", "")
	kidSecretName = env.GetString("KID_SECRET", "")

	os.Exit(m.Run())
}

func TestJWKSHandler(t *testing.T) {
	c := require.New(t)

	expectedJWKS := `{"keys":[{"kty":"RSA","kid":"123","use":"sig","n":"5l-M6MGnS6K8SNXUIqOGaaH_IO7NcBxwQJVd4X6uUcLHfdhyNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8_k4PoIUikiFnuCmkVDxcnl65_jv4DQtDL6GGqoLcYo2ENldfj8uDo09CmYS_DKuJxFyntaOREIMTaLQ3F72aDMk0ytVFu0cZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK_D8_ZKklkKTDnOmGlbVOKTvujH6fTJuQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqOWTuEpsvMNRubWoLt22c9f4PXaGDwqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98HQ","e":"AQAB"}]}`

	secretMock := secrets.NewSecretMock()
	secretMock.RegisterResponder(publicSecretName, func(ctx context.Context, name string) (string, error) {
		return "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5l+M6MGnS6K8SNXUIqOG\naaH/IO7NcBxwQJVd4X6uUcLHfdhyNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4\nPoIUikiFnuCmkVDxcnl65/jv4DQtDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFy\nntaOREIMTaLQ3F72aDMk0ytVFu0cZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8\n/ZKklkKTDnOmGlbVOKTvujH6fTJuQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqO\nWTuEpsvMNRubWoLt22c9f4PXaGDwqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98\nHQIDAQAB\n-----END PUBLIC KEY-----", nil
	})

	secretMock.RegisterResponder(privateSecretName, func(ctx context.Context, name string) (string, error) {
		return "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA5l+M6MGnS6K8SNXUIqOGaaH/IO7NcBxwQJVd4X6uUcLHfdhy\nNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4PoIUikiFnuCmkVDxcnl65/jv4DQt\nDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFyntaOREIMTaLQ3F72aDMk0ytVFu0c\nZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8/ZKklkKTDnOmGlbVOKTvujH6fTJu\nQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqOWTuEpsvMNRubWoLt22c9f4PXaGDw\nqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98HQIDAQABAoIBAQC8OffcuVVihC2I\n6UUxpCCPsG/PTa6HWoURD7msI6B0Z0wt86qkPCH0xlxufhxt/wk4GvIiEqm2P5YG\n7l0JUGh3EjuMHOMoQ+90rgkI+l7EqG3950OjtQvHP4aoF0BDlgZAv4h3FUl5dJsw\neow2mZfoVe/Zr4lz6YLze5ei3Z7J9YGjj62j7QGbKgbwPLqnqnNrUQqM4T0V9SaJ\nCE6sDxYo8M8kE2yqgiIvsA1D92u4AMchcdnjREBy/ogCRXzvuZQADC4st8UE4aFD\nmDNKwIbprymSa7atjSMz+lfWWBnuuzFmsf+72gXJVRRpmbm7onBxDHcJlqk0fMjv\nm+zuow4hAoGBAPyGYveulOLbwA+88DHBJ5GVKZNJHRGFsMWpysCC45OC5eKs9yrP\nnI6/0mvL1JFXvSbkkDql6qvlbKVcoH80h3ipFwuB4j8KbLXs/LiFMF0yg4IpAoq5\nGIp1RK6VBsv4OvP8vmkJxzakgRn5C7WyswNqJHGW108l7pmYEbOpfohZAoGBAOmL\nIE3ZbttXLOLEcHICcySpNdDTIWaiNqMtyjLLK29Ic8KDpXI1hfYzAw+7Q6y2QyEG\nl9l5IpkY6Wt29LMuWUMH6fnS1H/JeQOSnkT2y8PXSN95QKb2HtbP+ujqSdG12EPs\nILag0918ezcttGLszqWfipSZuSo2ZQ6b0A+uaxllAoGAECkBeFxBxurNNbSfom97\n+sMS8AwDwjVOBLhC82Ls8Wm1EHaFMsYqfLAl5SQcLFjzD+QcnsQzamC6PTLaSomw\nCba4dNIRCnu+TT4nRh+v4qby54d8VChYO7QZexqqXq86BpcsEEjB6OtKH8FiUHRp\nJFTMlEBU8wm4ZTfoGhlEsbECgYAkNJ1ddEfrWShsP2fvRNH07QaaySB0eNFfmsmt\n9jFVnzXTAfW0LvgFowLmfXGQZPEjPZJs9IqYkXQeZOKqpJTR/3gWcsjexq0sEJ7Y\nsioEwmtZucJ8H8vIIZYUZb3r9PUCEqk/ps8xlwrDEyLT80JWCtXBE9PQ533jNeSb\nib6wwQKBgQDkCKsHfxv/z+YgdMe3mUCSZi2gNttPQczjeUSYAxYITj/OJ1TfMuk4\n8gVdOcusHynFH3jEpnA8fqdpZpmhH/sAKPuQl/vwBCefVyBO5LkM14gxEIf9eq69\n7QVBd9ep1cN/5yYcJUJAcpjBxcbR8rXYowLtsYaGsC7G5tMlW8rJTg==\n-----END RSA PRIVATE KEY-----", nil
	})

	secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
		return "123", nil
	})

	logMock := logger.NewLoggerMock(nil)

	request := &requestJwksHandler{
		log:            logMock,
		secretsManager: secretMock,
	}

	ctx := context.Background()

	response, err := request.processJWKS(ctx, &apigateway.Request{})
	c.Nil(err)
	c.Equal(http.StatusOK, response.StatusCode)
	c.Equal(expectedJWKS, response.Body)
}

func TestJWKSHandlerFailed(t *testing.T) {
	c := require.New(t)

	secretMock := secrets.NewSecretMock()
	ctx := context.Background()
	logMock := logger.NewLoggerMock(nil)

	request := &requestJwksHandler{
		log:            logMock,
		secretsManager: secretMock,
	}

	t.Run("Public key fetching failed", func(t *testing.T) {
		dummyErr := fmt.Errorf("dummy error")

		secretMock.RegisterResponder(publicSecretName, func(ctx context.Context, name string) (string, error) {
			return "", dummyErr
		})

		defer secretMock.UnregisterResponder(publicSecretName)

		response, err := request.processJWKS(ctx, &apigateway.Request{})
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.NoError(err)
		c.Contains(logMock.Output.String(), "public_key_fetching_failed")
	})

	t.Run("Kid fetching failed", func(t *testing.T) {
		dummyErr := fmt.Errorf("dummy error")

		secretMock.RegisterResponder(publicSecretName, func(ctx context.Context, name string) (string, error) {
			return "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5l+M6MGnS6K8SNXUIqOG\naaH/IO7NcBxwQJVd4X6uUcLHfdhyNFNGEVFXodk9xhn0zJUxNtDzXlsw8aoC8/k4\nPoIUikiFnuCmkVDxcnl65/jv4DQtDL6GGqoLcYo2ENldfj8uDo09CmYS/DKuJxFy\nntaOREIMTaLQ3F72aDMk0ytVFu0cZ5Hyb24ixPBXhWHTMzsNG6yRO3uOVZqtK/D8\n/ZKklkKTDnOmGlbVOKTvujH6fTJuQ8T3p6jLI9J24K77fDlr6b38tZcDcKrhlAqO\nWTuEpsvMNRubWoLt22c9f4PXaGDwqRHo3SeBhb8YA0nSBEzNVgyt8iYfGq01tW98\nHQIDAQAB\n-----END PUBLIC KEY-----", nil
		})

		secretMock.RegisterResponder(kidSecretName, func(ctx context.Context, name string) (string, error) {
			return "", dummyErr
		})

		defer func() {
			secretMock.UnregisterResponder(publicSecretName)
			secretMock.UnregisterResponder(kidSecretName)
		}()

		response, err := request.processJWKS(ctx, &apigateway.Request{})
		c.Equal(http.StatusInternalServerError, response.StatusCode)
		c.NoError(err)
		c.Contains(logMock.Output.String(), "kid_fetching_failed")
	})
}
