FROM alpine:3.22.0

# hadolint ignore=DL3018
RUN apk add --no-cache git tini gnupg

COPY uplift_*.apk /tmp/
RUN apk add --no-cache --allow-untrusted /tmp/uplift_*.apk && \
    rm /tmp/uplift_*.apk

ENTRYPOINT ["/sbin/tini", "--", "uplift"]
CMD ["--help"]
