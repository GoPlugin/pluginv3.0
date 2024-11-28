# Docker environment

This folder contains a Plugin cluster environment created with [testcontainers-go](https://github.com/testcontainers/testcontainers-go/tree/main).

## CLI for Local Testing Environment

The command-line interface (CLI) located at `./integration-tests/docker/cmd/test_env.go` can be utilized to initiate a local testing environment. It is intended to replace Docker Compose in the near future.

Example:

```sh
# Set required envs
export PLUGIN_IMAGE="<plugin_node_docker_image_path>"
export PLUGIN_VERSION="<plugin_node_docker_image_version>" 
# Stream logs to Loki
export LOKI_TOKEN=...
export LOKI_URL=https://${loki_host}/loki/api/v1/push

cd ./integration-tests/docker/cmd

go run test_env.go start-env cl-cluster
```

## Obtaining Test Coverage for Plugin Node

To acquire test coverage data for end-to-end (E2E) tests on the Plugin Node, follow these steps:

1. Build Plugin Node docker image with the cover flag.

    First, build the Plugin Node Docker image with the `GO_COVER_FLAG` argument set to `true`. This enables the coverage flag in the build process. Here’s how you can do it:
    ```
    docker buildx build --platform linux/arm64 . -t localhost/plugin-local:develop -f ./core/plugin.Dockerfile --build-arg GO_COVER_FLAG=true
    ```
    Make sure to replace localhost/plugin-local:develop with the appropriate repository and tag.

2. Configure and Run E2E Tests
    Next, configure the E2E tests to generate an HTML coverage report. You need to modify the `overrides.toml` file as shown below to include the show_html_coverage_report setting under the `[Logging]` section:

    ```
    [Logging]
    show_html_coverage_report=true
    ```

After the tests are complete, the coverage report will be generated in HTML format. Example: `~/Downloads/go-coverage/TestOCRv2Basic_plugins-chain-reader/coverage.html`
```
    log.go:43: 16:29:46.73 INF Plugin node coverage html report saved filePath=~/Downloads/go-coverage/TestOCRv2Basic_plugins-chain-reader/coverage.html
```

