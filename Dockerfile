FROM alpine:3.16.1

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache \
    git=2.36.1-r0 \
    tini=0.19.0-r0

COPY uplift /usr/local/bin

ENTRYPOINT ["/sbin/tini", "--", "uplift"]
CMD ["--help"]