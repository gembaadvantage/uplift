# GPG Key fails to Import

Uplift supports the signing of commits by importing a GPG key and correctly configuring git. Your GPG key needs to be exported in the ASCII Armor Format (_optionally base64 encoded_) for this to work. Uplift will report the following error:

```text
uplift could not import GPG key with fingerprint FDA7347ACCE12A6CEBED57727B0EDBE188EE9114.
Check your GPG key was exported correctly.

For further details visit: https://upliftci.dev/faq/gpgimport
```

## How to fix it

You can resolve this error by exporting your key using the `--armor` flag. Please read the following [guide](../commit-signing.md#generating-a-gpg-key) on how to do this.
