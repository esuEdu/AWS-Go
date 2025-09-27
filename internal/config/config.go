package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Config struct {
	Addr              string
	DBHost            string
	DBPort            int
	DBUser            string
	DBName            string
	DBPassword        string
	DBPasswordSSMName string
	AWSRegion         string
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() (*Config, error) {
	c := &Config{}
	c.Addr = envOrDefault("ADDR", ":8080")
	c.DBHost = os.Getenv("DB_HOST")
	c.DBUser = envOrDefault("DB_USER", "appuser")
	c.DBName = envOrDefault("DB_NAME", "appdb")
	c.AWSRegion = os.Getenv("AWS_REGION")
	if p := os.Getenv("DB_PORT"); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			c.DBPort = i
		}
	}
	if c.DBPort == 0 {
		c.DBPort = 5432
	}
	c.DBPassword = os.Getenv("DB_PASSWORD")
	c.DBPasswordSSMName = os.Getenv("DB_PASSWORD_SSM_PARAM")

	if c.DBPassword == "" && c.DBPasswordSSMName != "" {
		pw, err := getSSMParameter(context.Background(), c.AWSRegion, c.DBPasswordSSMName)
		if err != nil {
			return nil, fmt.Errorf("get SSM param %s: %w", c.DBPasswordSSMName, err)
		}
		c.DBPassword = pw
	}

	if c.DBHost == "" || c.DBUser == "" || c.DBName == "" || c.DBPassword == "" {
		return nil, errors.New("missing required DB envs: DB_HOST, DB_USER, DB_NAME, DB_PASSWORD or DB_PASSWORD_SSM_PARAM")
	}
	if c.AWSRegion == "" {
		c.AWSRegion = envOrDefault("AWS_REGION", "us-east-1")
	}
	return c, nil
}

func getSSMParameter(ctx context.Context, region, name string) (string, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return "", err
	}
	client := ssm.NewFromConfig(cfg)
	out, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &name,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", errors.New("parameter has no value")
	}
	return *out.Parameter.Value, nil
}
