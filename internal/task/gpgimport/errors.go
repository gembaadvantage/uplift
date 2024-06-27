package gpgimport

import (
	"errors"
	"fmt"
)

// ErrGpgMissing is raised if git is not detected on the current $PATH
var ErrGpgMissing = errors.New("gpg is not currently installed under $PATH")

// ErrKeyImport is raised when a provided GPG key fails to be imported
type ErrKeyImport struct {
	fingerprint string
}

// Error returns a formatted message of the current error
func (e ErrKeyImport) Error() string {
	return fmt.Sprintf(`uplift could not import GPG key with fingerprint %s. Check your GPG
key was exported correctly.

For further details visit: https://upliftci.dev/faq/gpgimport
`, e.fingerprint)
}
