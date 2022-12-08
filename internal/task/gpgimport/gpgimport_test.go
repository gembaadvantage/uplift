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

package gpgimport

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/git"
	"github.com/gembaadvantage/uplift/internal/gpg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	assert.Equal(t, "importing gpg key", Task{}.String())
}

func TestSkip(t *testing.T) {
	t.Setenv("UPLIFT_GPG_KEY", "")
	t.Setenv("UPLIFT_GPG_PASSPHRASE", "")
	t.Setenv("UPLIFT_GPG_FINGERPRINT", "")

	assert.True(t, Task{}.Skip(&context.Context{}))
}

func TestSkipFalse(t *testing.T) {
	t.Setenv("UPLIFT_GPG_KEY", "key")
	t.Setenv("UPLIFT_GPG_PASSPHRASE", "passphrase")
	t.Setenv("UPLIFT_GPG_FINGERPRINT", "fingerprint")

	assert.False(t, Task{}.Skip(&context.Context{}))
}

func TestRun(t *testing.T) {
	git.InitRepo(t)

	t.Setenv("UPLIFT_GPG_KEY", gpg.TestKey)
	t.Setenv("UPLIFT_GPG_PASSPHRASE", gpg.TestPassphrase)
	t.Setenv("UPLIFT_GPG_FINGERPRINT", gpg.TestFingerprint)

	err := Task{}.Run(&context.Context{})

	require.NoError(t, err)
	assert.True(t, git.ConfigExists("user.signingKey", gpg.TestKeyID))
	assert.True(t, git.ConfigExists("commit.gpgsign", "true"))
	assert.True(t, git.ConfigExists("user.name", gpg.TestKeyUserName))
	assert.True(t, git.ConfigExists("user.email", gpg.TestKeyUserEmail))
}

func TestRunImportKeyFailed(t *testing.T) {
	git.InitRepo(t)

	t.Setenv("UPLIFT_GPG_KEY", "-----BEGIN PGP PRIVATE KEY BLOCK-----key-----END PGP PRIVATE KEY BLOCK-----")
	t.Setenv("UPLIFT_GPG_PASSPHRASE", "passphrase")
	t.Setenv("UPLIFT_GPG_FINGERPRINT", "AABBCCDDEEFF1122334455")

	err := Task{}.Run(&context.Context{})

	require.Error(t, err)
	assert.EqualError(t, err, `uplift could not import GPG key with fingerprint AABBCCDDEEFF1122334455. Check your GPG
key was exported correctly.

For further details visit: https://upliftci.dev/faq/gpgimport
`)
}

func TestRunGpgMissing(t *testing.T) {
	t.Setenv("PATH", "")

	err := Task{}.Run(&context.Context{})

	assert.EqualError(t, err, "gpg is not currently installed under $PATH")
}
