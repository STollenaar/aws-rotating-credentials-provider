package filecreds_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stollenaar/aws-rotating-credentials-provider/credentials/filecreds"
	"github.com/stretchr/testify/assert"
)

func TestFilecredentialsProvider_Retrieve(t *testing.T) {
	// Create a temporary directory for the test credentials file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test-credentials")

	// Write test credentials to the file
	testCredentials := `
[default]
aws_access_key_id = test-access-key-id
aws_secret_access_key = test-secret-access-key
aws_session_token = test-session-token
region = us-east-1
output = json
`
	err := os.WriteFile(testFilePath, []byte(testCredentials), 0644)
	assert.NoError(t, err, "Failed to write test credentials file")

	// Create a FilecredentialsProvider with the test file path
	provider := filecreds.NewFilecredentialsProvider(testFilePath)

	// Call Retrieve and validate the returned credentials
	creds, err := provider.Retrieve(context.Background())
	assert.NoError(t, err, "Retrieve should not return an error")
	assert.Equal(t, "test-access-key-id", creds.AccessKeyID, "AccessKeyID should match")
	assert.Equal(t, "test-secret-access-key", creds.SecretAccessKey, "SecretAccessKey should match")
	assert.Equal(t, "test-session-token", creds.SessionToken, "SessionToken should match")
	assert.Equal(t, filecreds.FilecredentialsName, creds.Source, "Source should match")
	assert.True(t, creds.CanExpire, "Credentials should be set to expire")
	assert.WithinDuration(t, time.Now().Add(2*time.Minute), creds.Expires, time.Second, "Expiration time should be approximately 2 minutes from now")
}

func TestFilecredentialsProvider_Retrieve_EmptyFilePath(t *testing.T) {
	// Create a FilecredentialsProvider with an empty file path
	provider := filecreds.NewFilecredentialsProvider("")

	// Call Retrieve and validate the returned error
	creds, err := provider.Retrieve(context.Background())
	assert.Error(t, err, "Retrieve should return an error for empty file path")
	assert.IsType(t, &filecreds.FilecredentialsEmptyError{}, err, "Error should be of type FilecredentialsEmptyError")
	assert.Equal(t, filecreds.FilecredentialsName, creds.Source, "Source should match")
}

func TestFilecredentialsProvider_Retrieve_FromRootFile(t *testing.T) {
	// Define the path to the test-credentials file in the root directory
	rootFilePath := filepath.Join("..", "..", "..", "test-credentials")

	// Ensure the file exists before running the test
	_, err := os.Stat(rootFilePath)
	assert.NoError(t, err, "test-credentials file should exist in the root directory")

	// Create a FilecredentialsProvider with the root file path
	provider := filecreds.NewFilecredentialsProvider(rootFilePath)

	// Call Retrieve and validate the returned credentials
	creds, err := provider.Retrieve(context.Background())
	assert.NoError(t, err, "Retrieve should not return an error")
	assert.Equal(t, "test-access-key-id", creds.AccessKeyID, "AccessKeyID should match")
	assert.Equal(t, "test-secret-access-key", creds.SecretAccessKey, "SecretAccessKey should match")
	assert.Equal(t, "test-session-token", creds.SessionToken, "SessionToken should match")
	assert.Equal(t, filecreds.FilecredentialsName, creds.Source, "Source should match")
	assert.True(t, creds.CanExpire, "Credentials should be set to expire")
	assert.WithinDuration(t, time.Now().Add(2*time.Minute), creds.Expires, time.Second, "Expiration time should be approximately 2 minutes from now")
}
