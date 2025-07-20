# go-sdk

Go SDK for the SECA API specification. This SDK is currently mostly handwritten but will be replaced using code generators for the majority of the code in the future. The client HTTP code is already generated using a code generator.

## Requirements

- **Go 1.23.1+** for basic functionality
- **Go 1.24.5+** required for spec generation with gomplate

### Go Version Compatibility

If you have an older Go version (< 1.24.5), you can use Go's automatic toolchain management to build the project:

```sh
GOTOOLCHAIN=auto make clean spec generate mock
```

This allows Go to automatically download and use the required Go version without changing your system installation.

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository:

    ```sh
    git clone git@github.com:eu-sovereign-cloud/go-sdk.git
    cd go-sdk
    ```

2. Initialize the submodule:

    ```sh
    git submodule init
    ```

3. Update all external dependencies:

    ```sh
    make update
    ```

4. Generate the API clients and mocks:

    ```sh
    make clean spec generate mock
    ```

## Testing

To execute unit and integration tests, run the following command:

```sh
make test
```
