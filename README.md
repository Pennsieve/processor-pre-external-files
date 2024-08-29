# External Files Pre-Processor

Retrieves one or more external files via HTTP and places it in the input directory. The external files must be specified
by a file called `config.json` in the config directory. We assume this file has been created by the workflow
manager or a previous processor.

To build:

`docker build -t pennsieve/external-files-pre-processor .`

On arm64 architectures:

`docker build -f Dockerfile_arm64 -t pennsieve/external-files-pre-processor .`

To run tests:

1. `cd service`
2. ` go test ./...`

To run integration test:

1. Copy `dev.env.example` to `dev.env`
2. Copy `config.json.example` to `config.json`
3. Run `./run-integration-test.sh dev.env`

The config in `config.json.example` contains several https://httpbin.org URLs which will download json
files that can be used to verify that query params and authentication are being handled correctly.

If you wish to test with other external files, then edit the contents of `config.json` and
run `./run-integration-test.sh dev.env` again.
