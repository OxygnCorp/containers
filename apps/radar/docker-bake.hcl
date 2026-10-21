target "docker-metadata-action" {}

variable "APP" {
  default = "radar"
}

variable "VERSION" {
  // renovate: datasource=github-releases depName=skyhook-io/radar
  default = "1.17.0"
}

variable "SOURCE" {
  default = "https://github.com/skyhook-io/radar"
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
  labels = {
    "org.opencontainers.image.source" = "${SOURCE}"
  }
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
  tags = ["${APP}:${VERSION}"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
  ]
}
