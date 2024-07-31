# External Files Pre-Processor

Retrieves one or more external files via HTTP and places it in the input directory.

To build:

`docker build -t pennsieve/external-files-pre-processor .`

On arm64 architectures:

`docker build -f Dockerfile_arm64 -t pennsieve/external-files-pre-processor .`

To run tests:

` go test ./...`

To run integration test:

1. Copy `dev.env.example` to `dev.env`
2. Run `./run-integration-test.sh dev.env`

The `EXTERNAL_FILES` value in `dev.env.example` contains several https://httpbin.org URLs which will download json
files that
can be used to verify that query params and authentication are being handled correctly.

If you wish to test with other external files, then edit the value of `EXTERNAL_FILES` in `dev.env` and
run `./run-integration-test.sh dev.env` again.
