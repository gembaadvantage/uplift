# Running with Docker

You can run Uplift directly from a docker image. Just mount your repository as a volume and set it as the working directory. 🐳

=== "DockerHub"

    ```sh
    docker run --rm -v $PWD:/tmp -w /tmp gembaadvantage/uplift release
    ```

=== "GHCR"

    ```sh
    docker run --rm -v $PWD:/tmp -w /tmp ghcr.io/gembaadvantage/uplift release
    ```

## Verifying with Cosign

Docker images can be verified using [cosign](https://github.com/sigstore/cosign).

=== "DockerHub"

    ```sh
    COSIGN_EXPERIMENTAL=1 cosign verify gembaadvantage/uplift
    ```

=== "GHCR"

    ```sh
    COSIGN_EXPERIMENTAL=1 cosign verify ghcr.io/gembaadvantage/uplift
    ```
