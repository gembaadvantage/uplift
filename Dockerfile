FROM alpine:3.17.3

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache \
    git=2.38.3-r1 \
    tini=0.19.0-r1 \
    gnupg=2.2.40-r0

COPY uplift_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/uplift_*.apk && \
    rm /tmp/uplift_*.apk

ENTRYPOINT ["/sbin/tini", "--", "uplift"]
CMD ["--help"]