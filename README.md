# External Files Pre-Processor

Retrieves one or more external files via HTTP and places it in the input directory.

To build:

`docker build -t pennsieve/external-files-pre-processor .`

On arm64 architectures:

`docker build -f Dockerfile_arm64 -t pennsieve/external-files-pre-processor .`

To run tests:

` go test ./...`

To run integration test:

1. Create an integration with the URLs you would like to download.
2. Copy `dev.env.example` to `dev.env`
3. In `dev.env` update `SESSION_TOKEN` with a valid token and `INTEGRATION_ID` with the id from the first step.
4. Run `./run-integration-test.sh dev.env`

