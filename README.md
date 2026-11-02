<div align="center">

## Containers

_OxygnCorp's opinionated collection of container images_

</div>

<div align="center">

![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/OxygnCorp/containers/release.yaml?style=for-the-badge&label=Release)

</div>

This repository re-packages the images that the [OxygnCorp fleet](https://github.com/OxygnCorp/home-ops) consumes but that upstream alone can't satisfy — an app graft, a missing distroless runtime or a multi-arch need. Each image is documented below: why it exists, what it grafts onto which source. Today:

| Image | What it does | Source |
| ----- | ------------ | ------ |
| [`radar`](https://github.com/OxygnCorp/containers/pkgs/container/radar) | [Radar](https://github.com/skyhook-io/radar) (distroless, no shell, no libc) with the OpenCode CLI — plus ripgrep and the musl loader/libraries they need — grafted in, so Radar's OSS **AI investigations** feature can spawn an agent inside the Radar pod | [`ghcr.io/skyhook-io/radar`](https://github.com/skyhook-io/radar) + [`ghcr.io/joryirving/opencode`](https://github.com/joryirving/home-ops) |

Browse the collection on the [GitHub Packages page for this repository's packages](https://github.com/orgs/OxygnCorp/packages?repo_name=containers).

The pipeline is maintained as a trimmed fork of [home-operations/containers](https://github.com/home-operations/containers): the build machinery (per-app `docker-bake.hcl`, bake metadata actions, release pipeline with attestations and SBOMs) is kept close to upstream, while apps, workflows and shared actions that make no sense here are pruned — every build targets a listed image.

## Features

### Tag Immutability

Containers built here do not use immutable tags in the traditional sense, as seen with [linuxserver.io](https://fleet.linuxserver.io/) or [Bitnami](https://bitnami.com/stacks/containers). Instead, consumers are expected to pin the `sha256` digest of the image. While this approach is less visually appealing, it ensures functionality and immutability.

| Container                                        | Immutable |
| ------------------------------------------------- | --------- |
| `ghcr.io/oxygncorp/radar:rolling`                 | ❌        |
| `ghcr.io/oxygncorp/radar:1.17.0`                  | ❌        |
| `ghcr.io/oxygncorp/radar:rolling@sha256:8053...`  | ✅        |
| `ghcr.io/oxygncorp/radar:1.17.0@sha256:8053...`   | ✅        |

_If pinning an image with a `sha256` digest, tools like [Renovate](https://github.com/renovatebot/renovate) — as configured in [`OxygnCorp/home-ops`](https://github.com/OxygnCorp/home-ops) — track the digest and the version together. For grafted images that expose the same upstream twice (a release tag and its registry digest), the release pipeline's renovate grouping rules in `.renovaterc.json5` keep them in lockstep so the digest never lags behind the version bump._

### Runtime identity

Images run as the user declared by the upstream image: the runtime identity, entrypoint and default command are re-declared explicitly in each graft so nothing is silently inherited. `radar` runs as `nonroot:nonroot` (uid 65532), like its distroless upstream:

```yaml
# Kubernetes, excerpt for the radar image
containers:
  - name: radar
    image: ghcr.io/oxygncorp/radar:1.17.0@sha256:<release-digest>
    securityContext:
      allowPrivilegeEscalation: false
      capabilities:
        drop: [ALL]
      readOnlyRootFilesystem: true
```

Passing arguments to the application follows the [Kubernetes documentation](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/):

```yaml
args:
  - --port
  - "9280"
```

### Verify Image Signature

These container images are attested using the [attest-build-provenance](https://github.com/actions/attest-build-provenance) action, with an SBOM alongside.

To verify that an image was built by this repository's CI, use the following command:

```sh
gh attestation verify --repo OxygnCorp/containers oci://ghcr.io/oxygncorp/${APP}:${TAG}
```

or by using [cosign](https://github.com/sigstore/cosign):

```sh
cosign verify-attestation --new-bundle-format --type slsaprovenance1 \
    --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
    --certificate-identity-regexp "^https://github.com/OxygnCorp/containers/.github/workflows/app-builder.yaml@refs/heads/main" \
    ghcr.io/oxygncorp/${APP}:${TAG}
```

## How it works

- **`apps/<name>/docker-bake.hcl`** declares the app's `VERSION` variable with a `renovate:` annotation; Renovate opens and auto-merges updates, which flow into the release pipeline on push (PR CI = lint + sandbox build + container tests; merge = tagged release build with attestations and SBOM).
- **`apps/<name>/Dockerfile`** pins upstream bases by digest and grafts what the app needs; Renovate tracks the digest pins. When a graft couples the same upstream twice (differently named release tag and registry digest), a **grouping rule** in `.renovaterc.json5` (e.g. `radar`) keeps the VERSION bump and the digest update in the same PR so the published tag always points at a manifest that exists.
- **`container_test.go`** lives beside each Dockerfile (shared helpers in `tests/`) and runs in the sandbox build on every PR.

### Adding an image

Consider adding an app when the upstream application is actively maintained **and** one of the following applies: no official image matches the fleet's constraints (distroless/static base, non-root by default, versioned tags), or the graft unlocks a fleet capability. The pipeline is app-agnostic: adding an app is just `apps/<name>/` (`Dockerfile`, `docker-bake.hcl`, `container_test.go`) plus a renovate grouping rule when the graft couples a version with a digest. An image that becomes unmaintained or unsustainable can be dropped from the listing above and from the build; published images remain available in the registry.

## Credits

This repository draws inspiration and ideas from [home-operations/containers](https://github.com/home-operations/containers), which produces the original form of the pipeline. Fleet-side conventions beyond the registry live in [`OxygnCorp/home-ops`](https://github.com/OxygnCorp/home-ops) (branch + PR flow, `flate`-validated GitOps, digest-pinned images in every consuming container).
