package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/esuEdu/aws-cognito-go/cognito"
)

func main() {
	action := flag.String("action", "", "Action to perform: signup | confirm | signin | refresh")
	username := flag.String("username", "", "Username (email)")
	password := flag.String("password", "", "Password")
	code := flag.String("code", "", "Confirmation code")
	refreshToken := flag.String("refresh-token", "", "Refresh token")

	flag.Parse()

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_APP_CLIENT_ID")

	if userPoolID == "" || clientID == "" {
		log.Fatal("Set COGNITO_USER_POOL_ID and COGNITO_APP_CLIENT_ID environment variables")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	cognitoClient := cognito.NewClient(cognitoidentityprovider.NewFromConfig(cfg), clientID, userPoolID)

	ctx := context.TODO()

	switch *action {
	case "signup":
		if err := cognitoClient.SignUp(ctx, *username, *password); err != nil {
			log.Fatalf("SignUp failed: %v", err)
		}
		fmt.Println("SignUp successful! Check your email for the verification code.")

	case "confirm":
		if err := cognitoClient.ConfirmSignUp(ctx, *username, *code); err != nil {
			log.Fatalf("ConfirmSignUp failed: %v", err)
		}
		fmt.Println("ConfirmSignUp successful! You can now sign in.")

	case "signin":
		tokens, err := cognitoClient.SignIn(ctx, *username, *password)
		if err != nil {
			log.Fatalf("SignIn failed: %v", err)
		}
		fmt.Println("SignIn successful!")
		fmt.Printf("Access Token: %s\n", *tokens.AccessToken)
		fmt.Printf("ID Token: %s\n", *tokens.IdToken)
		fmt.Printf("Refresh Token: %s\n", *tokens.RefreshToken)

	case "refresh":
		tokens, err := cognitoClient.RefreshToken(ctx, *refreshToken)
		if err != nil {
			log.Fatalf("Refresh failed: %v", err)
		}
		fmt.Println("Token refresh successful!")
		fmt.Printf("New Access Token: %s\n", *tokens.AccessToken)
		fmt.Printf("New ID Token: %s\n", *tokens.IdToken)

	default:
		log.Fatal("Invalid action. Use -action=signup | confirm | signin | refresh")
	}
}
