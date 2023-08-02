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
func (t Task) Skip(ctx *context.Context) bool {
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
