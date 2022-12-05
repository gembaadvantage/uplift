FROM alpine:3.17.0

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache \
    git=2.38.1-r0 \
    tini=0.19.0-r1 \
    gnupg=2.2.40-r0

COPY uplift /usr/local/bin

ENTRYPOINT ["/sbin/tini", "--", "uplift"]
CMD ["--help"]