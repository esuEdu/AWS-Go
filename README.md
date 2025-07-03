# Cognito CLI (Go)

A simple command-line interface in Go for managing AWS Cognito user authentication.

✅ Sign up users  
✅ Confirm accounts via verification code  
✅ Sign in and receive tokens  
✅ Refresh access tokens using refresh tokens

---

## ⚙️ AWS Setup (Required)

Before using the CLI, you must create a **Cognito User Pool** and **App Client**.

### 🔑 1. Configure AWS CLI

```bash
aws configure
Enter:

AWS Access Key ID

AWS Secret Access Key

Region (e.g. us-west-2)

Output format (e.g. json)
```

### 🧱 2. Create Cognito User Pool

```bash
aws cognito-idp create-user-pool \
  --pool-name my-cli-user-pool \
  --auto-verified-attributes email \
  --policies '{
    "PasswordPolicy": {
      "MinimumLength": 8,
      "RequireUppercase": true,
      "RequireLowercase": true,
      "RequireNumbers": true,
      "RequireSymbols": false
    }
  }' \
  --username-configuration '{"CaseSensitive": false}' \
  --query 'UserPool.Id' \
  --output text
```

#### Save the output (your User Pool ID):

```bash
export COGNITO_USER_POOL_ID="us-west-2_XXXXXXXXX"
```

### 🧩 3. Create Cognito App Client

```bash
aws cognito-idp create-user-pool-client \
  --user-pool-id $COGNITO_USER_POOL_ID \
  --client-name cli-client \
  --no-generate-secret \
  --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH \
  --query 'UserPoolClient.ClientId' \
  --output text
```

#### Save the output (your App Client ID):

```bash
export COGNITO_APP_CLIENT_ID="xxxxxxxxxxxxxxxxxxxxxx"
```

### 🌐 Set Environment Variables

```bash
export COGNITO_USER_POOL_ID="us-west-2_XXXXXXXXX"
export COGNITO_APP_CLIENT_ID="xxxxxxxxxxxxxxxxxxxxxx"
```

## 🚀 CLI Usage

### 📝 Sign Up

```bash
go run cmd/main.go -action=signup -username="email@example.com" -password="YourPass123!"
```

### ✅ Confirm Account

```bash
go run cmd/main.go -action=confirm -username="email@example.com" -code="123456"
```

### 🔐 Sign In

```bash
go run cmd/main.go -action=signin -username="email@example.com" -password="YourPass123!"
```

### 🔄 Refresh Tokens

```bash
go run cmd/main.go -action=refresh -refresh-token="your_refresh_token_here"
```

## 🛠 Build Binary

```bash
go build -o mycli cmd/main.go
```

#### Run it:

```bash
./mycli -action=signin -username="email@example.com" -password="YourPass123!"
```

## 📦 Flags Summary

Flag Description
-action signup, confirm, signin, refresh
-username Email or username
-password Password (or Refresh Token for refresh)
-code Confirmation code
-refresh-token Refresh token (used in refresh action)

## 🧼 Notes

Passwords must meet Cognito’s policy.

App client should not use a client secret.

Make sure the app client has these auth flows enabled:

ALLOW_USER_PASSWORD_AUTH

ALLOW_REFRESH_TOKEN_AUTH
