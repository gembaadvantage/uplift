/*
Copyright (c) 2022 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gpg

import (
	"os"
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

func TestDeleteKey(t *testing.T) {
	importKey(t)

	err := DeleteKey(TestFingerprint)
	require.NoError(t, err)

	_, err = Clean(Run("--batch", "--list-secret-keys", TestFingerprint))
	assert.EqualError(t, err, "gpg: error reading key: No secret key")
}

func importKey(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	fi, err := os.CreateTemp(dir, "key.asc")
	require.NoError(t, err)

	err = os.WriteFile(fi.Name(), []byte(TestKey), 0o600)
	require.NoError(t, err)

	_, err = Run("--batch", "--import", "--yes", fi.Name())
	require.NoError(t, err)
}
