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
