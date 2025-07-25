package gpg

import (
	"encoding/base64"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	secLineRegex = regexp.MustCompile("(?im)^sec.*")
	uidLineRegex = regexp.MustCompile("(?im)^uid.*")
	uidRegex     = regexp.MustCompile(`([^\(]*)(\s\(.*\)\s)?<(.*)>`)

	importKeyPath      = filepath.Join(os.TempDir(), "uplift-gpg-import.asc")
	activateKeyPath    = filepath.Join(os.TempDir(), "uplift-activate-key.txt")
	activateKeySigPath = filepath.Join(os.TempDir(), "uplift-activate-key.txt.sig")
)

// KeyDetails contains details about an imported private key
type KeyDetails struct {
	ID        string
	UserName  string
	UserEmail string
}

// IsInstalled identifies whether gpg is installed under the current $PATH
func IsInstalled() bool {
	_, err := Run("--version")
	return err == nil
}

// ImportKey attempts to import the private key using the provided passphrase.
// If importing is successful, the key will automatically be activated ready
// for use
func ImportKey(key, passphrase, fingerprint string) (KeyDetails, error) {
	if !strings.HasPrefix(key, "--") {
		// Decode base64 string into expected format
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return KeyDetails{}, err
		}
		key = string(decoded)
	}

	if err := os.WriteFile(importKeyPath, []byte(key), 0o600); err != nil {
		return KeyDetails{}, err
	}
	defer os.Remove(importKeyPath)

	os.TempDir()

	// Import the key using the temporary file on disk
	if _, err := Run("--batch", "--import", "--yes", importKeyPath); err != nil {
		return KeyDetails{}, err
	}

	out, _ := Clean(Run("--batch", "--with-colons", "--list-secret-keys", fingerprint))

	// Parse the key ID and the user details from the GPG private key
	sec := secLineRegex.FindString(out)
	uid := uidLineRegex.FindString(out)
	uid = strings.Split(uid, ":")[9]

	uidParts := uidRegex.FindStringSubmatch(uid)

	details := KeyDetails{
		ID:        strings.Split(sec, ":")[4],
		UserName:  strings.TrimSpace(uidParts[1]),
		UserEmail: uidParts[3],
	}

	// Activate the newly imported key
	if err := os.WriteFile(activateKeyPath, []byte(`hello, world!`), 0o600); err != nil {
		return details, err
	}
	defer os.Remove(activateKeyPath)

	_, err := Run("--local-user",
		fingerprint,
		"--batch",
		"--detach-sig",
		"--pinentry-mode",
		"loopback",
		"--no-tty",
		"--passphrase",
		passphrase,
		activateKeyPath)
	if err != nil {
		return details, err
	}
	defer os.Remove(activateKeySigPath)

	return details, nil
}

// Run executes a gpg command and returns its output or errors
func Run(args ...string) (string, error) {
	cmd := exec.Command("gpg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// RunAgent executes a gpg-agent command and returns its output or errors
func RunAgent(args ...string) (string, error) {
	cmd := exec.Command("gpg-agent", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// Clean the output
func Clean(output string, err error) (string, error) {
	// Preserve multi-line output, but trim the trailing newline
	output = strings.TrimSuffix(strings.ReplaceAll(output, "'", ""), "\n")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}
