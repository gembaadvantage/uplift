package gpg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInstalled(t *testing.T) {
	assert.True(t, IsInstalled())
}

func TestIsInstalled_NotOnPath(t *testing.T) {
	t.Setenv("PATH", "")

	assert.False(t, IsInstalled())
}

func TestImportKey(t *testing.T) {
	t.Cleanup(func() {
		Run("--batch", "--yes", "--delete-secret-keys", TestFingerprint)
		Run("--batch", "--yes", "--delete-keys", TestFingerprint)
	})

	details, err := ImportKey(TestKey, TestPassphrase, TestFingerprint)

	require.NoError(t, err)
	assert.Equal(t, "AAC7E54CBD73F690", details.ID)
	assert.Equal(t, "john.smith", details.UserName)
	assert.Equal(t, "john.smith@testing.com", details.UserEmail)
}

func TestImportKeyBase64(t *testing.T) {
	t.Cleanup(func() {
		Run("--batch", "--yes", "--delete-secret-keys", TestFingerprint)
		Run("--batch", "--yes", "--delete-keys", TestFingerprint)
	})

	details, err := ImportKey(TestKeyBase64, TestPassphrase, TestFingerprint)

	require.NoError(t, err)
	assert.Equal(t, TestKeyID, details.ID)
	assert.Equal(t, TestKeyUserName, details.UserName)
	assert.Equal(t, TestKeyUserEmail, details.UserEmail)
}
