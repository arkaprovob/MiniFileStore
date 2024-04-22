# MiniStore Utility Application

MiniStore is a small file store utility application designed to store files along with their related meta-information. It provides users with a command-line interface (CLI) to interact with the application via HTTP requests.
## Usage

1. **Build the Image Locally:**
    - Clone the repository or unzip the project.
    - Navigate to the project's root directory and run the following command:
      ```bash
      podman build -t ministore .
      ```

2. **Running the Locally Built Container:**
    - After building the image, execute the following command:
      ```bash
      podman run -d -p 8080:8080 -v </path/to/local/store/files>:/home/appuser/store/files:z -v </path/to/local/store/record>:/home/appuser/store/record:z ministore
      ```
    - Replace `</path/to/local/store/files>` and `</path/to/local/store/record>` with the actual paths to the folders on your machine. These folders will be mounted within the container.

3. **Using the Pre-built Container:**
    - A public container image is available at `quay.io/arbhatta/minifs:latest`.
    - To use this image, execute the following command:
      ```bash
      podman run --name minifscont -d -p 8080:8080 -v </path/to/local/store/files>:/home/appuser/store/files:z -v </path/to/local/store/record>:/home/appuser/store/record:z quay.io/arbhatta/minifs:latest
      ```
    - Replace `</path/to/local/store/files>` and `</path/to/local/store/record>` with the actual paths on your machine.

4. **Dealing with Permission Denied Issues in the Container:**
    - If you encounter permission denied issues in the container logs during any operation, it might be due to the container not being able to access the mounted volume.
    - In such cases, stop the container and execute the following commands:
      ```bash
      podman unshare chown 1000:1000 </path/to/local/store/files>
      podman unshare chown 1000:1000 </path/to/local/store/record>
      ```
    - This will change the ownership of the mounted volume to the user in the container. Then start the container again.

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