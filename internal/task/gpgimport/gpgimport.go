package gpgimport

import (
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/gpg"
)

const (
	envGpgKey         = "UPLIFT_GPG_KEY"
	envGpgPassphrase  = "UPLIFT_GPG_PASSPHRASE"
	envGpgFingerprint = "UPLIFT_GPG_FINGERPRINT"
)

// Task for tagging a repository
type Task struct{}

// String generates a string representation of the task
func (t Task) String() string {
	return "importing gpg key"
}

// Skip running the task if no version has changed
func (t Task) Skip(_ *context.Context) bool {
	return strings.TrimSpace(os.Getenv(envGpgKey)) == "" ||
		strings.TrimSpace(os.Getenv(envGpgPassphrase)) == "" ||
		strings.TrimSpace(os.Getenv(envGpgFingerprint)) == ""
}

// Run the task to import a provided gpg key and enable gpg signing of commits
func (t Task) Run(ctx *context.Context) error {
	log.Debug("checking if gpg is installed")
	if !gpg.IsInstalled() {
		return ErrGpgMissing
	}

	key := os.Getenv(envGpgKey)
	passphrase := os.Getenv(envGpgPassphrase)
	fingerprint := os.Getenv(envGpgFingerprint)

	log.WithField("fingerprint", fingerprint).Info("importing gpg key")
	keyDetails, err := gpg.ImportKey(key, passphrase, os.Getenv(envGpgFingerprint))
	if err != nil {
		return ErrKeyImport{fingerprint: fingerprint}
	}

	log.Info("setting git config to enable gpg signing")
	return ctx.GitClient.ConfigSetL("user.signingKey", keyDetails.ID,
		"commit.gpgsign", "true",
		"user.name", keyDetails.UserName,
		"user.email", keyDetails.UserEmail)
}
