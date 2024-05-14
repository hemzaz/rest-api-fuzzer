# REST API Fuzzer

## Overview

The REST API Fuzzer is a CLI tool designed to probe and test REST APIs for structure discovery and security vulnerabilities. The tool operates in two modes: `probe` and `breach`. The `probe` mode maps the API structure and types, while the `breach` mode attempts to breach the security of the API using the data collected during the probe.

## Features

- **Probe Mode**: Discover endpoints, HTTP methods, and data types of the REST API.
- **Breach Mode**: Test the security of the API by sending malicious payloads to discovered endpoints.
- **Multi-threaded Execution**: Utilize concurrent threads to speed up testing.
- **AI and Machine Learning**: Generate sophisticated payloads based on previous responses and known attack patterns.
- **Persistent Data Storage**: Store probe results and breach attempts using BoltDB.
- **HTTPS Support**: Supports both HTTP and HTTPS URLs.
- **Authentication Support**: Supports Basic Auth, Bearer Tokens, and more.
- **Rate Limiting and Throttling**: Configurable request rate to avoid overwhelming the target API.
- **Custom Headers and Parameters**: Allows specifying custom headers and parameters for requests.
- **Detailed Logging**: Logs errors and results for easier debugging and analysis.

## Installation

### Prerequisites

- Go 1.16 or later
- Git

### Steps

1. Clone the repository:
    ```
    git clone https://github.com/yourusername/rest-api-fuzzer.git
    cd rest-api-fuzzer
    ```

2. Initialize and download dependencies:
    ```
    make init
    ```

3. Build the project:
    ```
    make build
    ```

## Usage

### General Usage

```
./fuzzer [command] [flags]
```

### Commands

#### Probe

Probes the structure of a REST API by sending requests to common endpoints and recording the responses.

```
./fuzzer probe --url <base_url>
```

**Flags:**
- `--url`, `-u`: Base URL of the API (required).
- `--auth-type`, `-a`: Authentication type (e.g., 'Bearer', 'Basic').
- `--auth-token`, `-k`: Authentication token.
- `--rate-limit`, `-r`: Rate limit in requests per second.
- `--header`: Custom headers (key:value).

**Example:**
```
./fuzzer probe --url https://localhost:8080 --auth-type Bearer --auth-token mytoken --rate-limit 10 --header "X-Custom-Header: value"
```

**How it works:**
1. The `probe` mode sends HTTP requests to the provided base URL using common HTTP methods (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD).
2. It records the responses to determine which endpoints and methods are supported.
3. The results, including the URL, method, response status, and response body, are stored in a BoltDB database for later use.

#### Breach

Attempts to breach the security of a REST API by sending malicious payloads to discovered endpoints.

```
./fuzzer breach --url <base_url> --threads <number_of_threads>
```

**Flags:**
- `--url`, `-u`: Base URL of the API (required).
- `--threads`, `-t`: Number of concurrent threads (default is 4).
- `--auth-type`, `-a`: Authentication type (e.g., 'Bearer', 'Basic').
- `--auth-token`, `-k`: Authentication token.
- `--header`: Custom headers (key:value).

**Example:**
```
./fuzzer breach --url https://localhost:8080 --threads 8 --auth-type Bearer --auth-token mytoken --header "X-Custom-Header: value"
```

**How it works:**
1. The `breach` mode retrieves the endpoints and methods discovered during the `probe` phase from the BoltDB database.
2. It generates malicious payloads using AI and machine learning techniques, incorporating mutation-based fuzzing and pattern-based payload generation.
3. These payloads are then sent to the discovered endpoints using multiple concurrent threads.
4. The responses are analyzed to check for potential security vulnerabilities such as SQL injection, XSS, and privilege escalation.
5. The results, including the payload, response status, and response body, are stored in the BoltDB database.

### Makefile Commands

- **Initialize the project**: This will initialize the Go module and get the necessary packages.
    ```
    make init
    ```

- **Build the project**: This will compile the Go code and generate the binary.
    ```
    make build
    ```

- **Run the Probe Mode**: This will execute the fuzzer in probe mode to discover the API structure.
    ```
    make probe
    ```

- **Run the Breach Mode**: This will execute the fuzzer in breach mode to test the API's security.
    ```
    make breach
    ```

- **Clean Up**: This will remove the generated binary and the database file.
    ```
    make clean
    ```

- **Format the Code**: This will format the Go code.
    ```
    make fmt
    ```

- **Run Tests**: This will run any tests in the project.
    ```
    make test
    ```

- **Complete Workflow**: This will run the complete workflow from initialization to execution.
    ```
    make full
    ```

## How It Works

### Probe Mode

1. **Sending Requests**: The fuzzer sends HTTP requests to the provided base URL using common methods (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD).
2. **Recording Responses**: It records the responses to determine which endpoints and methods are supported.
3. **Storing Results**: The results, including the URL, method, response status, and response body, are stored in a BoltDB database for later use in the breach phase.

### Breach Mode

1. **Retrieving Probe Data**: The fuzzer retrieves the endpoints and methods discovered during the probe phase from the BoltDB database.
2. **Generating Malicious Payloads**: It generates sophisticated malicious payloads using AI and machine learning techniques, such as mutation-based fuzzing and pattern-based payload generation.
3. **Sending Payloads**: These payloads are then sent to the discovered endpoints using multiple concurrent threads.
4. **Analyzing Responses**: The responses are analyzed to check for potential security vulnerabilities such as SQL injection, XSS, and privilege escalation.
5. **Storing Results**: The results, including the payload, response status, and response body, are stored in the BoltDB database.


## Contribution

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
