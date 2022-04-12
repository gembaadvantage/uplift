FROM alpine:3.15.0

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache \
    git=2.34.2-r0 \
    tini=0.19.0-r0

COPY uplift /

ENTRYPOINT ["/sbin/tini", "--", "/uplift"]
CMD ["--help"]