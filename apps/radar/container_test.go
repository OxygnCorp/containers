package main

import (
	"testing"

	helpers "github.com/home-operations/containers/tests"
)

func Test(t *testing.T) {
	image := helpers.GetTestImage("ghcr.io/oxygncorp/radar:rolling")

	// Go binaries and grafted files must all be present; the Radar image is
	// distroless/static (no shell), so everything is asserted through image
	// filesystem inspection rather than being invoked.
	helpers.RequireFileExists(t, image, "/radar")
	helpers.RequireFileExists(t, image, "/usr/local/bin/opencode")
	helpers.RequireFileExists(t, image, "/usr/bin/rg")

	// musl loader and its libraries, copied from the Alpine OpenCode stage.
	helpers.RequireFileExists(t, image, "/lib/ld-musl-x86_64.so.1")
	helpers.RequireFileExists(t, image, "/usr/lib/libstdc++.so.6")
	helpers.RequireFileExists(t, image, "/usr/lib/libgcc_s.so.1")
	helpers.RequireFileExists(t, image, "/usr/lib/libpcre2-8.so.0")

	// ripgrep is the one grafted binary that runs standalone (no HOME, no
	// network, no writes): an entrypoint run exercises the musl loader and
	// the shared libraries end to end.
	helpers.RequireCommandSucceeds(t, image, nil, "/usr/bin/rg", "--version")
}
