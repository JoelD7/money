package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	c := require.New(t)

	event := events.APIGatewayCustomAuthorizerRequest{
		Type:               "",
		AuthorizationToken: "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6InJlYWQgd3JpdGUiLCJpc3MiOiJodHRwczovLzM4cXNscGU4ZDkuZXhlY3V0ZS1hcGkudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vc3RhZ2luZyIsInN1YiI6InRlc3RAZ21haWwuY29tIiwiYXVkIjoiaHR0cHM6Ly9sb2NhbGhvc3Q6MzAwMCIsImV4cCI6MTcwODExNzY5NywibmJmIjoxNjc3MDE1NDk3LCJpYXQiOjE2NzcwMTM2OTd9.QOIY5X_ruqwX5N9GSFqlgp7YcY6GylrPYz7Z9XaozXKjaU0_sRm56P1yM0nhz01LN1nXiyoKxIk915_5i8wcaz_WnTGqM3fu51eXocl-9-X1uLqTl5y8oE1u6F-aA3Us_BkTY_5UrVV1TnKqndKliOV86ZobFSeJnwpHfaeukxOXYDkU7c7xcI-1ZLqFazReNRR83-I9TO7ZWsgRK-blgJYN-nBtmw9_EzDISVOsk7k7zVfszuRdA-10WAVBlPjQwvATkxEzXyhc3aYnUAFBG-Xrz8SWIAbhIi4uxJj5EvufqmbyhRO3yrTnZUygH_fT_331uZ_ryT6563nbhhVaeQ",
		MethodArn:          "",
	}

	_, err := handleRequest(context.Background(), event)
	c.Nil(err)
}
