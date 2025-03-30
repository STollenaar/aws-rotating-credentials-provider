package filecreds

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/viper"
)

const (
	// FilecredentialsName provides a name of Static provider
	FilecredentialsName = "Filecredentials"
)

type DefaultCredentials struct {
	Default Credentials `mapstructure:"default"`
}

type Credentials struct {
	AccessKeyID     string `mapstructure:"aws_access_key_id"`
	SecretAccessKey string `mapstructure:"aws_secret_access_key"`
	SessionToken    string `mapstructure:"aws_session_token"`
	Region          string `mapstructure:"region"`
	Output          string `mapstructure:"output"`
}

// FilecredentialsEmptyError is emitted when static credentials are empty.
type FilecredentialsEmptyError struct{}

func (*FilecredentialsEmptyError) Error() string {
	return "rotating credentials are empty"
}

// A FilecredentialsProvider is a set of credentials which are set, and will
// never expire.
type FilecredentialsProvider struct {
	FilePath string
}

// NewFilecredentialsProvider return a FilecredentialsProvider initialized with the AWS
// credentials passed in.
func NewFilecredentialsProvider(file string) FilecredentialsProvider {
	return FilecredentialsProvider{
		FilePath: file,
	}
}

// Retrieve returns the credentials or error if the credentials are invalid.
func (s FilecredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	// fmt.Println("Fetching Credentials")
	if s.FilePath == "" {
		return aws.Credentials{
			Source: FilecredentialsName,
		}, &FilecredentialsEmptyError{}
	}

	var creds DefaultCredentials
	viper.SetConfigName(path.Base(s.FilePath))
	viper.AddConfigPath(path.Dir(s.FilePath))
	viper.SetConfigType("ini")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error reading creds: %e\n", err)
		return aws.Credentials{
			Source: FilecredentialsName,
		}, err
	}

	err = viper.Unmarshal(&creds)
	if err != nil {
		fmt.Printf("error unmarshalling creds: %e\n", err)
		return aws.Credentials{
			Source: FilecredentialsName,
		}, err
	}

	newT := time.Now().Add(time.Duration(time.Minute * 2))
	// fmt.Printf("Credentials will expire at: %v\n", newT.Format(time.RFC1123))
	return aws.Credentials{
		AccessKeyID:     creds.Default.AccessKeyID,
		SecretAccessKey: creds.Default.SecretAccessKey,
		SessionToken:    creds.Default.SessionToken,
		Source:          FilecredentialsName,

		CanExpire: true,
		Expires:   newT,
	}, nil
}
