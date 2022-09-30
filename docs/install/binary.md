# Installing the Binary

You can use various package managers to install the Uplift binary. Take your pick.

## Package Managers

### Homebrew

To use [Homebrew](https://brew.sh/):

```sh
brew install gembaadvantage/tap/uplift
```

### Scoop

To use [Scoop](https://scoop.sh/):

```sh
scoop install uplift
```

### Apt

To install using the [apt](https://ubuntu.com/server/docs/package-management) package manager:

```sh
echo 'deb [trusted=yes] https://fury.upliftci.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/uplift.list
sudo apt update
sudo apt install -y uplift
```

You may need to install the `ca-certificates` package if you encounter [trust issues](https://gemfury.com/help/could-not-verify-ssl-certificate/) with regards to the Gemfury certificate:

```sh
sudo apt update && sudo apt install -y ca-certificates
```

### Yum

To install using the yum package manager:

```sh
echo '[uplift]
name=uplift
baseurl=https://fury.upliftci.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/uplift.repo
sudo yum install -y uplift
```

### Aur

To install from the [aur](https://archlinux.org/) using [yay](https://github.com/Jguer/yay):

```sh
yay -S uplift-bin
```

### Linux Packages

Download and manually install one of the `.deb`, `.rpm` or `.apk` packages from the [Releases](https://github.com/gembaadvantage/uplift/releases) page.

=== "Apt"

    ```sh
    sudo apt install uplift_*.deb
    ```

=== "Yum"

    ```sh
    sudo yum localinstall uplift_*.rpm
    ```

=== "Apk"

    ```sh
    sudo apk add --no-cache --allow-untrusted uplift_*.apk
    ```

### Bash Script

To install the latest version using a bash script:

```sh
curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash
```

Download a specific version using the `-v` flag. The script uses `sudo` by default but can be disabled through the `--no-sudo` flag.

```sh
curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash -s -- -v v2.6.3 --no-sudo
```

## Manual Download of Binary

Binary downloads of uplift can be found on the [Releases](https://github.com/gembaadvantage/uplift/releases) page. Unpack the uplift binary and add it to your `PATH`.

## Verifying a Binary with Cosign

All binaries can be verified using [cosign](https://github.com/sigstore/cosign).

1. Download the checksum files that need to be verified:
   ```sh
   curl -sL https://github.com/gembaadvantage/uplift/releases/download/v2.5.0/checksums.txt -O
   curl -sL https://github.com/gembaadvantage/uplift/releases/download/v2.5.0/checksums.txt.sig -O
   curl -sL https://github.com/gembaadvantage/uplift/releases/download/v2.5.0/checksums.txt.pem -O
   ```

1. Verify the signature of the checksum file:
   ```sh
   cosign verify-blob --cert checksums.txt.pem --signature checksums.txt.sig checksums.txt
   ```

1. Download any release artefact and verify its SHA256 signature matches the entry within the checksum file:
   ```sh
   sha256sum --ignore-missing -c checksums.txt
   ```
