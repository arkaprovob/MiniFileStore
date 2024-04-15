# MiniStore Utility Application

MiniStore is a small file store utility application designed to store files along with their related meta-information. It provides users with a command-line interface (CLI) to interact with the application via HTTP requests.
## Usage

To build a Docker image of MiniStore, clone the repository or unzip the project, and then run the following command on the project's root.

```bash
podman build -t ministore .
```

To run the Docker image, execute the following command:

```bash
podman run -d -p 8080:8080 \
  -v /path/to/local/store:/app/store/files \
  -v /path/to/local/record:/app/store/record \
  ministore
```

Replace `/path/to/local/store` and `/path/to/local/record` with the actual paths to the folders on your machine. These folders will be mounted within the Docker container.

A public container image is available at `quay.io/arbhatta/minifs:latest`. To use this image, execute the following command:
```bash
podman run -d -p 8080:8080 \
  -v </path/to/local/store>:/app/store/files \
  -v </path/to/local/record>:/app/store/record \
  quay.io/arbhatta/minifs:latest

```
This command will download the `quay.io/arbhatta/minifs:latest` image from Quay and start a container from it. The `-v` options are for mounting your local directories into the container. Replace `/path/to/local/store` and `/path/to/local/record` with the actual paths on your machine.
## API Routes

MiniStore exposes the following API routes:

- `/`: Root endpoint. Accessing this endpoint provides information about the application.
- `/api/v1/store`: Handle storing files along with their meta-information.
- `/api/v1/update`: Update existing files in the store with new content or meta-information.
- `/api/v1/exists`: Check the existence of a file in the store.
- `/api/v1/list`: List all files stored in the application.
- `/api/v1/delete`: Delete a file from the store.
- `/api/v1/frequency`: Calculate the frequency of words in the stored files.

All API details are available in `api-specs.yaml` in the form of OpenAPI v3.0.0 specifications. To access the API specifications, simply navigate to the root path (`/`) of the running Docker/Podman instance. For example, if MiniStore is running on `localhost` and port `8080`, you can access the API specs by visiting `http://localhost:8080/`.

## Scope of Improvement

- Code Refactoring: Refactor the codebase to improve readability, maintainability, and scalability.
- Reduce Redundancy: Identify and eliminate redundant code to streamline the application.
- Increase Test Coverage: Add more test cases to ensure comprehensive test coverage and reliability.
- Improve Unit Tests: Enhance the maturity of unit test cases to validate individual components more effectively.