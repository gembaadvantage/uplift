FROM alpine:3.15.0

# Install tini to ensure docker waits for uplift to finish before terminating
RUN apk add --no-cache git tini

ENTRYPOINT ["/sbin/tini", "--", "uplift"]
CMD ["--help"]