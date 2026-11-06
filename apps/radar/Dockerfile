# syntax=docker/dockerfile:1

# Radar with embedded OpenCode for the OSS "AI investigations" feature.
#
# Upstream ghcr.io/skyhook-io/radar is a distroless/static image (no shell,
# no libc). Radar v1.15.0+ spawns a supported agent CLI found on its PATH
# at startup; this graft copies the `opencode` CLI (and `rg`, which
# OpenCode shells out to) plus the musl loader and shared libraries they
# need, extracted from the Alpine-based upstream OpenCode image built by
# @joryirving. If a new musl dependency shows up at runtime (e.g. after an
# OpenCode upgrade — verify with `ldd` on the binaries), add a COPY here.
FROM ghcr.io/joryirving/opencode:2.2.5@sha256:0bba6c502886d834d6358f71d36e82cbd9ac5cd4e54252bed0734a60e2b86b02 AS opencode

FROM ghcr.io/skyhook-io/radar:1.19.0@sha256:e8dc4d847916f66833f6ec9b18aaeb426257063572693c3c8866ec43f8bca7dc AS radar

FROM radar

EXPOSE 9280

COPY --from=opencode /usr/local/bin/opencode /usr/local/bin/opencode
COPY --from=opencode /usr/bin/rg /usr/bin/rg
COPY --from=opencode /lib/ld-musl-x86_64.so.1 /lib/ld-musl-x86_64.so.1
COPY --from=opencode /usr/lib/libstdc++.so.6 /usr/lib/libstdc++.so.6
COPY --from=opencode /usr/lib/libgcc_s.so.1 /usr/lib/libgcc_s.so.1
COPY --from=opencode /usr/lib/libpcre2-8.so.0 /usr/lib/libpcre2-8.so.0

USER nonroot:nonroot
# Kept identical to the upstream radar image so the Helm chart keeps its
# expected runtime shape; the digest-pinned FROM already provides it, but
# being explicit removes the runtime-identity drift risk.
ENTRYPOINT ["/radar", "--listen-address=0.0.0.0"]
CMD ["--no-browser"]

# test build for libc dev headers fix
