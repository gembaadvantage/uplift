package gpgimport

import (
	"testing"

	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/gpg"
	"github.com/purpleclay/gitz/gittest"
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
	gittest.InitRepository(t)

	t.Setenv("UPLIFT_GPG_KEY", gpg.TestKey)
	t.Setenv("UPLIFT_GPG_PASSPHRASE", gpg.TestPassphrase)
	t.Setenv("UPLIFT_GPG_FINGERPRINT", gpg.TestFingerprint)

	err := Task{}.Run(&context.Context{})

	require.NoError(t, err)
	gittest.MustExec(t, "")
	assert.Equal(t, gpg.TestKeyID, gittest.MustExec(t, "git config --get user.signingKey"))
	assert.Equal(t, "true", gittest.MustExec(t, "git config --get commit.gpgsign"))
	assert.Equal(t, gpg.TestKeyUserName, gittest.MustExec(t, "git config --get user.name"))
	assert.Equal(t, gpg.TestKeyUserEmail, gittest.MustExec(t, "git config --get user.email"))
}

func TestRunImportKeyFailed(t *testing.T) {
	gittest.InitRepository(t)

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
