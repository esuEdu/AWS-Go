package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type Client struct {
	svc        *cognitoidentityprovider.Client
	clientID   string
	userPoolID string
}

func NewClient(svc *cognitoidentityprovider.Client, clientID, userPoolID string) *Client {
	return &Client{
		svc:        svc,
		clientID:   clientID,
		userPoolID: userPoolID,
	}
}

// SignUp registers a new user with email/username and password
func (c *Client) SignUp(ctx context.Context, username, password string) error {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []cognitoTypes.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(username),
			},
		},
	}

	_, err := c.svc.SignUp(ctx, input)
	if err != nil {
		return fmt.Errorf("sign up failed: %w", err)
	}

	return nil
}

// ConfirmSignUp confirms user registration with the verification code
func (c *Client) ConfirmSignUp(ctx context.Context, username, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(code),
	}

	_, err := c.svc.ConfirmSignUp(ctx, input)
	if err != nil {
		return fmt.Errorf("confirm sign up failed: %w", err)
	}

	return nil
}

// SignIn authenticates a user and returns tokens
func (c *Client) SignIn(ctx context.Context, username, password string) (*cognitoTypes.AuthenticationResultType, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	resp, err := c.svc.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("sign in failed: %w", err)
	}

	return resp.AuthenticationResult, nil
}

// RefreshToken uses a refresh token to get new access/id tokens
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*cognitoTypes.AuthenticationResultType, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
		},
	}

	resp, err := c.svc.InitiateAuth(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("refresh token failed: %w", err)
	}

	return resp.AuthenticationResult, nil
}
