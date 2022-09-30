# Install

There are many different ways to install uplift. You can install the binary using either a supported package manager, manually, or by compiling the source yourself.

## Installing the binary

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

You may need to install the `ca-certificates` package if you encounter [trust issues](https://gemfury.com/help/could-not-verify-ssl-certificate/) with regards to the gemfury certificate:

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

```sh
sudo apt install uplift_*.deb
```

```sh
sudo yum localinstall uplift_*.rpm
```

```sh
sudo apk add --no-cache --allow-untrusted uplift_*.apk
```

### Bash Script

To install the latest version using a bash script:

```sh
curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash
```

A specific version can be downloaded by using the `-v` flag. By default the script uses `sudo`, which can be turned off by using the `--no-sudo` flag.

```sh
curl https://raw.githubusercontent.com/gembaadvantage/uplift/main/scripts/install | bash -s -- -v v2.6.3 --no-sudo
```

### Manually

Binary downloads of uplift can be found on the [Releases](https://github.com/gembaadvantage/uplift/releases) page. Unpack the uplift binary and add it to your `PATH`.

## Compiling from source

Uplift is written using [Go 1.18+](https://go.dev/doc/install) and should be installed along with [go-task](https://taskfile.dev/#/installation), as it is preferred over using make.

Then clone the code from GitHub:

```sh
git clone https://github.com/gembaadvantage/uplift.git
cd uplift
```

Build uplift:

```sh
task
```

And check that everything works:

```sh
./bin/uplift version
```

!!!tip "Fancy Contributing?"

    Since you have the code checked out and locally built, you are only one step away from contributing. Take a peek at the [Contributing Guide](https://github.com/gembaadvantage/uplift/blob/main/CONTRIBUTING.md)

## Verifying Artefacts

All verification is carried out using cosign and it must be [installed](https://docs.sigstore.dev/cosign/installation) before proceeding.

### Binaries

All binaries can be verified using the checksum file, which has been signed using cosign.

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

!!!tip "Don't mix versions"

    For checksum verification to work, all artefacts must be downloaded from the same release

### Docker

Docker images can be verified using cosign directly, as the signature will be embedded within the docker manifest.

!!!info "Cosign Verification"

    Cosign verification was introduced to all docker images from version `v2.5.0`

=== "DockerHub"
    ```sh
    COSIGN_EXPERIMENTAL=1 cosign verify gembaadvantage/uplift
    ```

=== "GHCR"
    ```sh
    COSIGN_EXPERIMENTAL=1 cosign verify ghcr.io/gembaadvantage/uplift
    ```

## Running with Docker

You can run uplift directly from a docker image. Depending on how you have cloned the repository, you may need to tweak the following command to work for your setup.

=== "DockerHub"
    ```sh
    docker run --rm -v $PWD:/tmp -w /tmp gembaadvantage/uplift release
    ```

=== "GHCR"
    ```sh
    docker run --rm -v $PWD:/tmp -w /tmp ghcr.io/gembaadvantage/uplift release
    ```

!!!warning "Issue with SSH Cloned Repositories"

    Outstanding issue with pushing changes back to a cloned SSH repository, see: [#148](https://github.com/gembaadvantage/uplift/issues/148)

## Oh My Zsh

Install the custom uplift [plugin](https://github.com/gembaadvantage/uplift-oh-my-zsh) for full autocompletion support.
