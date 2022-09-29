FROM alpine:3.16.2

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache \
    git=2.36.2-r0 \
    tini=0.19.0-r0 \
    gnupg=2.2.35-r4 \
    bash=5.1.16-r2

COPY uplift /usr/local/bin

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD ["--help"]